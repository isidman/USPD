package mesh

import (
	"crypto/rand"
	"encoding/binary"
	"sync"

	"meshfeed/internal/domain"
	"meshfeed/internal/store"
)

// DefaultTTL bounds how many hops a packet can travel before nodes stop
// relaying it — Meshtastic calls this a hop limit for the same reason:
// without it, a flood-routed mesh never quiets down.
const DefaultTTL = 3

// seenCacheCapacity is how many recent message IDs a node remembers for
// dedup. 512 IDs at 4 bytes each is 2KB — deliberately small enough to be
// plausible on a microcontroller, not just a demo laptop.
const seenCacheCapacity = 512

// Node is one participant in the mesh: it can originate a feed, relay
// packets it doesn't care about (extending the mesh's reach for everyone),
// and subscribe to feeds it wants to actively pull. There is no separate
// "subscription cursor" data structure — a feed's entry in the local
// FeedStore, specifically its LatestSeq, *is* the cursor. That avoids two
// pieces of state that could drift out of sync.
type Node struct {
	ID        string
	Store     *store.FeedStore
	transport Transport

	mu            sync.Mutex
	subscriptions map[domain.FeedID]bool
	origins       map[domain.FeedID]bool

	seen *seenCache
	Log  func(format string, args ...any)
}

func NewNode(id string, transport Transport) *Node {
	return &Node{
		ID:            id,
		Store:         store.New(),
		transport:     transport,
		subscriptions: make(map[domain.FeedID]bool),
		origins:       make(map[domain.FeedID]bool),
		seen:          newSeenCache(seenCacheCapacity),
		Log:           func(string, ...any) {},
	}
}

// Subscribe marks a feed as one this node actively wants: on receiving an
// Advertise for it that's ahead of what's held locally, the node will send
// a Want rather than just relaying the advertisement onward.
func (n *Node) Subscribe(feedID domain.FeedID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.subscriptions[feedID] = true
}

func (n *Node) isSubscribed(feedID domain.FeedID) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.subscriptions[feedID]
}

// Publish appends a new item to a feed this node originates and returns
// it. It does not push the item over the air — that happens on the next
// AdvertiseAll, when a subscriber (if any is in range, now or later via a
// relay) learns the new LatestSeq and pulls it. This mirrors RSS: the
// publisher doesn't push to every reader; readers notice something's new
// and fetch it.
func (n *Node) Publish(feedID domain.FeedID, content []byte, timestamp uint32) (domain.Item, error) {
	n.mu.Lock()
	n.origins[feedID] = true
	n.mu.Unlock()

	item := domain.Item{
		Seq:       n.Store.LatestSeq(feedID) + 1,
		Timestamp: timestamp,
		Content:   content,
	}
	n.Store.Append(feedID, item)
	return item, nil
}

// AdvertiseAll broadcasts one Advertise packet per feed this node holds
// any items for — origin feeds and cached/relayed ones alike. Advertising
// cached feeds too is what lets any node act as a store-and-forward point,
// the same role Meshtastic's Store & Forward module plays.
func (n *Node) AdvertiseAll() error {
	for _, feedID := range n.Store.Feeds() {
		pkt := domain.Packet{
			Type:      domain.PacketAdvertise,
			FeedID:    feedID,
			TTL:       DefaultTTL,
			MsgID:     randomMsgID(),
			LatestSeq: n.Store.LatestSeq(feedID),
		}
		n.seen.SeenBefore(pkt.MsgID) // don't reprocess our own broadcast if it echoes back
		if err := n.broadcastPacket(pkt); err != nil {
			return err
		}
	}
	return nil
}

// Run reads from the transport's inbox until it's closed, handling each
// packet as it arrives. Call it in a goroutine; it returns when the inbox
// channel closes.
func (n *Node) Run() {
	for data := range n.transport.Inbox() {
		if err := n.HandleIncoming(data); err != nil {
			n.Log("meshfeed[%s]: dropping malformed packet: %v", n.ID, err)
		}
	}
}

// HandleIncoming processes one packet received from the transport. It's
// exported directly (not just reachable via Run) so tests can drive it
// deterministically without a goroutine and a channel race.
func (n *Node) HandleIncoming(data []byte) error {
	pkt, err := domain.Unmarshal(data)
	if err != nil {
		return err
	}
	if n.seen.SeenBefore(pkt.MsgID) {
		return nil // already processed (or it's our own broadcast echoing back)
	}

	switch pkt.Type {
	case domain.PacketAdvertise:
		n.handleAdvertise(pkt)
	case domain.PacketWant:
		n.handleWant(pkt)
	case domain.PacketItem:
		n.handleItem(pkt)
	}

	return n.maybeRelay(pkt)
}

func (n *Node) handleAdvertise(pkt domain.Packet) {
	if !n.isSubscribed(pkt.FeedID) {
		return
	}
	local := n.Store.LatestSeq(pkt.FeedID)
	if pkt.LatestSeq <= local {
		return // already caught up
	}
	want := domain.Packet{
		Type:    domain.PacketWant,
		FeedID:  pkt.FeedID,
		TTL:     DefaultTTL,
		MsgID:   randomMsgID(),
		FromSeq: local,
	}
	n.seen.SeenBefore(want.MsgID)
	n.Log("meshfeed[%s]: behind on %s (have %d, saw %d) — sending want", n.ID, pkt.FeedID, local, pkt.LatestSeq)
	_ = n.broadcastPacket(want)
}

func (n *Node) handleWant(pkt domain.Packet) {
	items := n.Store.ItemsSince(pkt.FeedID, pkt.FromSeq)
	for _, item := range items {
		resp := domain.Packet{
			Type:   domain.PacketItem,
			FeedID: pkt.FeedID,
			TTL:    DefaultTTL,
			MsgID:  randomMsgID(),
			Item:   item,
		}
		n.seen.SeenBefore(resp.MsgID)
		_ = n.broadcastPacket(resp)
	}
}

func (n *Node) handleItem(pkt domain.Packet) {
	if n.Store.Append(pkt.FeedID, pkt.Item) {
		n.Log("meshfeed[%s]: received %s#%d: %q", n.ID, pkt.FeedID, pkt.Item.Seq, pkt.Item.Content)
	}
}

// maybeRelay rebroadcasts a packet with a decremented TTL so it keeps
// propagating outward — this is what makes it a *mesh*: a node relays
// traffic for feeds it doesn't itself care about, extending everyone
// else's reach. TTL 0 stops the flood from running forever.
func (n *Node) maybeRelay(pkt domain.Packet) error {
	if pkt.TTL == 0 {
		return nil
	}
	pkt.TTL--
	return n.broadcastPacket(pkt)
}

func (n *Node) broadcastPacket(pkt domain.Packet) error {
	data, err := pkt.Marshal()
	if err != nil {
		return err
	}
	return n.transport.Broadcast(data)
}

func randomMsgID() uint32 {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return binary.BigEndian.Uint32(b[:])
}

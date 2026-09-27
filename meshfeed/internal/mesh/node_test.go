package mesh_test

import (
	"testing"

	"meshfeed/internal/domain"
	"meshfeed/internal/mesh"
	"meshfeed/internal/simtransport"
)

// pump drains every node's inbox and processes each packet, round after
// round, until nothing moved in a full pass — simulating packets
// propagating (and relaying) through the mesh until it settles.
func pump(t *testing.T, nodes map[string]*mesh.Node, transports map[string]*simtransport.Transport) {
	t.Helper()
	const maxRounds = 20
	for round := 0; round < maxRounds; round++ {
		progressed := false
		for id, node := range nodes {
			tr := transports[id]
		drain:
			for {
				select {
				case data := <-tr.Inbox():
					progressed = true
					if err := node.HandleIncoming(data); err != nil {
						t.Fatalf("node %s: HandleIncoming: %v", id, err)
					}
				default:
					break drain
				}
			}
		}
		if !progressed {
			return
		}
	}
	t.Fatalf("pump: did not settle within %d rounds (possible relay loop)", maxRounds)
}

func TestDirectSync_OnlyHappensInRange(t *testing.T) {
	net := simtransport.NewNetwork()
	publisherT := net.NewTransport("publisher")
	subscriberT := net.NewTransport("subscriber")

	publisher := mesh.NewNode("publisher", publisherT)
	subscriber := mesh.NewNode("subscriber", subscriberT)

	feed := domain.NewFeedID("garden-club")
	subscriber.Subscribe(feed)

	if _, err := publisher.Publish(feed, []byte("compost bin needs turning"), 1_700_000_000); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	nodes := map[string]*mesh.Node{"publisher": publisher, "subscriber": subscriber}
	transports := map[string]*simtransport.Transport{"publisher": publisherT, "subscriber": subscriberT}

	// Out of range: advertising shouldn't reach the subscriber at all.
	if err := publisher.AdvertiseAll(); err != nil {
		t.Fatalf("AdvertiseAll: %v", err)
	}
	pump(t, nodes, transports)
	if got := subscriber.Store.LatestSeq(feed); got != 0 {
		t.Fatalf("expected subscriber to have nothing while out of range, got latest seq %d", got)
	}

	// Move into range and advertise again — now it should sync automatically.
	net.SetInRange("publisher", "subscriber", true)
	if err := publisher.AdvertiseAll(); err != nil {
		t.Fatalf("AdvertiseAll: %v", err)
	}
	pump(t, nodes, transports)

	items := subscriber.Store.ItemsSince(feed, 0)
	if len(items) != 1 || string(items[0].Content) != "compost bin needs turning" {
		t.Fatalf("expected subscriber to have synced the item, got %+v", items)
	}
}

func TestRelay_ReachesSubscriberThroughIntermediateNode(t *testing.T) {
	net := simtransport.NewNetwork()
	publisherT := net.NewTransport("publisher")
	relayT := net.NewTransport("relay")
	subscriberT := net.NewTransport("subscriber")

	publisher := mesh.NewNode("publisher", publisherT)
	relay := mesh.NewNode("relay", relayT)
	subscriber := mesh.NewNode("subscriber", subscriberT)

	// A chain: publisher <-> relay <-> subscriber. Publisher and subscriber
	// are never directly in range of each other — only the relay can hear
	// both. This is the actual "mesh" property: the relay doesn't
	// subscribe to anything, it just forwards.
	net.SetInRange("publisher", "relay", true)
	net.SetInRange("relay", "subscriber", true)

	feed := domain.NewFeedID("watershed-alerts")
	subscriber.Subscribe(feed)

	if _, err := publisher.Publish(feed, []byte("water level rising past the footbridge"), 1_700_000_500); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	nodes := map[string]*mesh.Node{"publisher": publisher, "relay": relay, "subscriber": subscriber}
	transports := map[string]*simtransport.Transport{"publisher": publisherT, "relay": relayT, "subscriber": subscriberT}

	if err := publisher.AdvertiseAll(); err != nil {
		t.Fatalf("AdvertiseAll: %v", err)
	}
	pump(t, nodes, transports)

	items := subscriber.Store.ItemsSince(feed, 0)
	if len(items) != 1 || string(items[0].Content) != "water level rising past the footbridge" {
		t.Fatalf("expected subscriber to receive the item via the relay, got %+v", items)
	}

	// The relay itself never subscribed — it should still have cached the
	// item purely by forwarding it, which is what lets it act as a
	// store-and-forward point for a node that comes into range later.
	relayed := relay.Store.ItemsSince(feed, 0)
	if len(relayed) != 1 {
		t.Fatalf("expected relay to have cached the item while forwarding it, got %+v", relayed)
	}
}

func TestUnsubscribedFeed_IsNotPulled(t *testing.T) {
	net := simtransport.NewNetwork()
	publisherT := net.NewTransport("publisher")
	bystanderT := net.NewTransport("bystander")
	net.SetInRange("publisher", "bystander", true)

	publisher := mesh.NewNode("publisher", publisherT)
	bystander := mesh.NewNode("bystander", bystanderT)

	feed := domain.NewFeedID("unrelated-feed")
	// bystander deliberately does not subscribe.
	if _, err := publisher.Publish(feed, []byte("hello"), 1_700_000_000); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	nodes := map[string]*mesh.Node{"publisher": publisher, "bystander": bystander}
	transports := map[string]*simtransport.Transport{"publisher": publisherT, "bystander": bystanderT}

	if err := publisher.AdvertiseAll(); err != nil {
		t.Fatalf("AdvertiseAll: %v", err)
	}
	pump(t, nodes, transports)

	if got := bystander.Store.LatestSeq(feed); got != 0 {
		t.Fatalf("expected bystander not to pull an unsubscribed feed, got latest seq %d", got)
	}
}

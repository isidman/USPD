// Package store holds items a node has locally — whether it originated
// them, or cached them while relaying or subscribing. Keyed by sequence
// number rather than a plain slice so items that arrive out of order
// (normal on a lossy multi-hop mesh) don't need reordering logic here.
package store

import (
	"sort"
	"sync"

	"meshfeed/internal/domain"
)

type FeedStore struct {
	mu    sync.Mutex
	items map[domain.FeedID]map[uint32]domain.Item
}

func New() *FeedStore {
	return &FeedStore{items: make(map[domain.FeedID]map[uint32]domain.Item)}
}

// Append stores an item and reports whether it was new — a duplicate
// (same feed, same seq, already held) reports false so callers (the mesh
// engine) know not to re-relay it.
func (s *FeedStore) Append(feedID domain.FeedID, item domain.Item) (isNew bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	feed, ok := s.items[feedID]
	if !ok {
		feed = make(map[uint32]domain.Item)
		s.items[feedID] = feed
	}
	if _, exists := feed[item.Seq]; exists {
		return false
	}
	feed[item.Seq] = item
	return true
}

// LatestSeq is the highest sequence number held for a feed, or 0 if none.
func (s *FeedStore) LatestSeq(feedID domain.FeedID) uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var max uint32
	for seq := range s.items[feedID] {
		if seq > max {
			max = seq
		}
	}
	return max
}

// ItemsSince returns every held item with Seq > fromSeq, ascending.
func (s *FeedStore) ItemsSince(feedID domain.FeedID, fromSeq uint32) []domain.Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []domain.Item
	for seq, item := range s.items[feedID] {
		if seq > fromSeq {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Seq < out[j].Seq })
	return out
}

// Feeds lists every feed this node holds at least one item for.
func (s *FeedStore) Feeds() []domain.FeedID {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.FeedID, 0, len(s.items))
	for id := range s.items {
		out = append(out, id)
	}
	return out
}

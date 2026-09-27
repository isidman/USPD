package mesh

import "sync"

// seenCache remembers recently-seen message IDs so a node relays each
// packet at most once — the "listen before repeating" half of Meshtastic's
// managed flood routing. Bounded, not a map that grows forever: a real
// node has kilobytes of RAM, not gigabytes.
type seenCache struct {
	mu       sync.Mutex
	order    []uint32
	set      map[uint32]struct{}
	capacity int
}

func newSeenCache(capacity int) *seenCache {
	return &seenCache{set: make(map[uint32]struct{}, capacity), capacity: capacity}
}

// SeenBefore reports whether id was already recorded, and records it if
// not. One call does both jobs so callers can't check-then-forget-to-mark
// between two separate calls.
func (c *seenCache) SeenBefore(id uint32) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.set[id]; ok {
		return true
	}
	c.set[id] = struct{}{}
	c.order = append(c.order, id)
	if len(c.order) > c.capacity {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.set, oldest)
	}
	return false
}

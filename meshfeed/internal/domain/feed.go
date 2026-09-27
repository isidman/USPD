// Package domain holds the types that mean the same thing on every node,
// independent of radio hardware, storage, or transport — same separation
// toolshed uses, for the same reason (CLAUDE.md rule 4: modular).
package domain

import (
	"crypto/sha256"
	"encoding/hex"
)

// FeedID is a short fingerprint, not a human-readable name — LoRa's payload
// budget (see internal/domain/packet.go) can't afford to spend bytes on
// strings that fit in a wire packet only once. A node that knows a feed's
// name maps it to this ID locally; only the ID ever goes over the air.
type FeedID [8]byte

func NewFeedID(name string) FeedID {
	sum := sha256.Sum256([]byte(name))
	var id FeedID
	copy(id[:], sum[:8])
	return id
}

func (id FeedID) String() string {
	return hex.EncodeToString(id[:])
}

// Item is one entry in a feed's append-only log — the unit an "RSS reader"
// would call an article. Timestamp is coarse (unix seconds, not
// nanoseconds) because those bytes are also coming out of the same tiny
// packet budget.
type Item struct {
	Seq       uint32
	Timestamp uint32
	Content   []byte
}

// Feed is local bookkeeping: the human-readable name behind a FeedID, and
// whether this node is the one originating it (vs. just relaying or
// subscribing to it).
type Feed struct {
	ID     FeedID
	Name   string
	Origin bool
}

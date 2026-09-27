// Package mesh is the protocol engine: advertise what you have, pull what
// you're missing, relay what passes through — independent of whether
// packets travel over a real LoRa radio or a simulated one in tests.
package mesh

// Transport is how packets actually move. A real implementation wraps a
// LoRa radio driver (e.g. talking to an SX127x/RFM95 module over SPI or a
// serial AT-command firmware) — not built here, see BLUEPRINT.md. The
// simulated implementation in internal/simtransport is what this
// package's own tests, and cmd/demo, run against.
type Transport interface {
	// Broadcast sends raw bytes to whatever's in range right now — LoRa is
	// a shared-medium broadcast radio, there's no "send to one peer."
	Broadcast(data []byte) error
	// Inbox delivers packets as they're received. The engine reads from
	// it in a loop for as long as the node runs.
	Inbox() <-chan []byte
}

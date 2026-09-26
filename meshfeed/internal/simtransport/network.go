// Package simtransport simulates LoRa radio range for tests and the
// cmd/demo walkthrough — no real hardware is available in the environment
// this blueprint was written in, so the mesh.Transport interface is
// implemented against a controllable in-memory "who can hear whom" model
// instead. A real implementation would wrap a LoRa radio driver; nothing
// in internal/mesh needs to change to swap it in later (same seam
// toolshed used for storage — see meshfeed/BLUEPRINT.md).
package simtransport

import "sync"

// Network is the shared simulated airwave. Nodes start out of range of
// everyone by default — call SetInRange to move them into (or out of)
// radio contact, the way walking two LoRa devices toward each other would.
type Network struct {
	mu      sync.Mutex
	nodes   map[string]*Transport
	inRange map[[2]string]bool
}

func NewNetwork() *Network {
	return &Network{
		nodes:   make(map[string]*Transport),
		inRange: make(map[[2]string]bool),
	}
}

// NewTransport registers a node on this simulated network and returns its
// mesh.Transport implementation.
func (n *Network) NewTransport(nodeID string) *Transport {
	t := &Transport{id: nodeID, network: n, inbox: make(chan []byte, 256)}
	n.mu.Lock()
	n.nodes[nodeID] = t
	n.mu.Unlock()
	return t
}

// SetInRange controls whether two nodes can currently hear each other.
// Symmetric, like real radio range.
func (n *Network) SetInRange(a, b string, inRange bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.inRange[rangeKey(a, b)] = inRange
}

func rangeKey(a, b string) [2]string {
	if a < b {
		return [2]string{a, b}
	}
	return [2]string{b, a}
}

// Transport is one node's view of the simulated network. It implements
// mesh.Transport.
type Transport struct {
	id      string
	network *Network
	inbox   chan []byte
}

func (t *Transport) Broadcast(data []byte) error {
	t.network.mu.Lock()
	defer t.network.mu.Unlock()
	for id, other := range t.network.nodes {
		if id == t.id {
			continue
		}
		if !t.network.inRange[rangeKey(t.id, id)] {
			continue
		}
		select {
		case other.inbox <- data:
		default:
			// Inbox full: simulates real radio congestion / packet loss
			// rather than blocking the sender forever.
		}
	}
	return nil
}

func (t *Transport) Inbox() <-chan []byte {
	return t.inbox
}

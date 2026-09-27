# meshfeed: an RSS feed for a LoRa mesh

## The idea

Subscribe to a node. When you come into radio range of it — or of anyone
relaying its updates — you get whatever you missed, automatically. No
internet, no server, no manual sync button. That's the whole product.

It's RSS's model (a feed is an append-only sequence of items; a reader
tracks how far it's read; catching up means fetching what's new) carried
onto a radio that can move maybe a few hundred bytes at a time and isn't
always in earshot of the thing you're subscribed to.

## Prior art this borrows from, on purpose

Two existing systems already solve most of this, in different halves:

- **[Secure Scuttlebutt](https://en.wikipedia.org/wiki/Secure_Scuttlebutt)**
  is almost exactly the data model: each identity owns an append-only,
  cryptographically signed feed; peers gossip and replicate by comparing
  how far along a feed they've each seen, and syncing the gap when they
  meet (over LAN, Bluetooth, or a USB stick — "sneakernet"). meshfeed's
  `Item`/sequence-number model is this, unsigned for now (see "What's not
  built here").
- **[Meshtastic](https://meshtastic.org/docs/overview/mesh-algo/)** is the
  real, running proof that "managed flood routing" — every node listens
  briefly, relays once if nobody beat it to it, and stops after a hop
  limit — works on LoRa at real-world mesh sizes (100+ nodes). Its
  [Store & Forward module](https://meshtastic.org/docs/configuration/module/store-and-forward-module/)
  — a node that caches messages so a device reconnecting later can catch
  up — is exactly what falls out of meshfeed's design for free: any node
  that relays a feed ends up holding a cache of it, whether or not it
  subscribed.

meshfeed is those two ideas combined: SSB's per-feed pull-sync, carried
over Meshtastic's flood-relay mechanism, on a wire format small enough for
LoRa.

## The physical constraint that shapes everything

LoRa's usable payload is roughly 51–222 bytes depending on spreading
factor and region, and in the EU868 band a 1% duty-cycle limit caps how
often a node can transmit *at all*, independent of nominal bitrate — see
[The Things Network's regional parameters](https://www.thethingsnetwork.org/docs/lorawan/regional-parameters/eu868/)
for the exact numbers. That rules out JSON (the per-field key overhead
alone would blow a 200-byte budget) and rules out "just re-send the whole
feed" sync strategies. Everything in `internal/domain/packet.go` is a
hand-rolled binary format sized against `MaxPacketBytes = 200`, with a
14-byte shared header (type, feed ID, hop count, dedup ID) leaving the
rest for whatever the packet type actually needs.

## What's actually built here

A complete, tested vertical slice of the sync-and-relay mechanic, in Go,
plus a minimal TypeScript reader UI:

- `internal/domain` — `FeedID` (an 8-byte fingerprint, not a name — names
  cost bytes a wire packet can't spare), `Item`, and the three packet
  types (`Advertise`, `Want`, `Item`) with hand-written `Marshal`/
  `Unmarshal` against the LoRa payload budget.
- `internal/store` — an in-memory, append-only, per-feed item cache, keyed
  by sequence number so out-of-order arrival (normal on a lossy multi-hop
  mesh) doesn't need special handling.
- `internal/mesh` — the actual protocol engine (`Node`): advertise what
  you hold, send a `Want` when an advertisement shows you're behind on a
  feed you subscribed to, answer `Want`s for feeds you hold, and relay
  every packet once (decrementing a hop count, deduping by message ID) —
  regardless of whether you yourself care about that feed. That last part
  is what makes it a *mesh* and not just point-to-point sync: a node with
  zero subscriptions still extends everyone else's reach.
- `internal/simtransport` — a controllable simulated radio: mark any two
  nodes in or out of range and watch sync behavior change accordingly.
  This is the same swappable-seam pattern toolshed used for storage
  (CLAUDE.md rule 5: isolate what's likely to break) — a real LoRa
  `Transport` implementation drops in without `internal/mesh` changing at
  all.
- `internal/api` + `web/` — a small **read-only** HTTP view and a
  framework-free TypeScript reader UI onto a node's local feed store: the
  actual "RSS reader" experience. Read-only on purpose — subscribing is a
  node-configuration decision (who do I want to hear from), not a web
  form.
- `cmd/demo` — a runnable, three-node scenario proving the specific
  behavior asked for: a subscriber gets nothing while out of range, then
  syncs automatically through an intermediate relay node that never
  subscribed to anything, once in range — then serves the result over
  HTTP so the reader UI can show it.

Tests cover: packet round-trip encode/decode and the payload-budget
boundary, direct sync only happening in range, multi-hop relay through an
uninvolved intermediate node, an unsubscribed feed correctly *not* being
pulled, and the HTTP API. All of it — `go test ./...`, `npm run build`,
`vitest run`, and a live run of `cmd/demo` piped into the actual reader UI
— was run for real; see the PR this shipped in for the output.

## What's not built here

- **Cryptographic signing.** SSB signs every message so a feed can't be
  forged by whoever relays it. This slice trusts content at face value —
  fine for a demo, not fine for anything where a malicious relay matters.
  Ed25519 (stdlib `crypto/ed25519`) is the natural fit: sign on `Publish`,
  verify on `handleItem` before accepting into the store.
- **Duty-cycle-aware pacing.** `handleWant` fires off every missing item
  back-to-back. A real EU868 deployment needs to spread those
  transmissions out to respect the 1% duty cycle — a rate limiter between
  `Node` and `Transport.Broadcast`, not a change to the protocol logic.
- **A real LoRa `Transport`.** Implement the two-method interface in
  `internal/mesh/transport.go` against a radio driver (RadioLib over
  serial to an SX127x/RFM95 module is the common hobbyist path). Nothing
  else in this repository needs to change.
- **Feed discovery.** Right now a node has to be told a `FeedID` (derived
  from a name it already knows) to subscribe to it. Advertising
  human-readable names, or a directory-of-feeds mechanism, is a real
  design question — worth its own writeup rather than bolted on here
  speculatively (CLAUDE.md rule 6).

## Adding a new packet type

1. Add the constant to `domain.PacketType` in `internal/domain/packet.go`.
2. Add its encode/decode branches to `Packet.Marshal` and `Unmarshal`,
   sized against `MaxPacketBytes`.
3. Add a `case` to the switch in `mesh.Node.HandleIncoming`. The relay
   step (`maybeRelay`) already applies to every packet type automatically
   — you don't need to touch it.

## Running it

See `README.md` in this directory.

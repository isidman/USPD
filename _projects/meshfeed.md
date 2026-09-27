---
title: "meshfeed"
tagline: "An RSS feed for a LoRa mesh: subscribe to a node, get its updates automatically when in range — directly or through a relay."
category: tools
status: idea
license: "AGPL-3.0-or-later (code), CC BY-SA 4.0 (docs)"
repo_url: "https://github.com/isidman/USPD/tree/main/meshfeed"
docs_url: "https://github.com/isidman/USPD/blob/main/meshfeed/BLUEPRINT.md"
tech_stack:
  - "Go (backend/protocol engine, zero third-party dependencies)"
  - "TypeScript (read-only reader UI, no framework)"
principles:
  repairable: true
  common_components: true
  parameterized: true
  modular: true
  isolates_failure: true
  minimal: true
  learnable: true
  legible: true
---

LoRa mesh radios can carry a message a few kilometers with no internet and
no infrastructure — but "receive updates from a node automatically" isn't
solved by the radio alone; it needs a sync protocol suited to tiny packets,
scarce airtime, and nodes that aren't always in range of each other.

meshfeed combines two proven ideas rather than inventing from scratch:
[Secure Scuttlebutt](https://en.wikipedia.org/wiki/Secure_Scuttlebutt)'s
model (each feed is an append-only log; peers gossip and pull-sync the gap
when they meet) carried over
[Meshtastic](https://meshtastic.org/docs/overview/mesh-algo/)'s proven
flood-relay mechanism (every node relays once, with a hop limit and a
dedup cache), on a hand-rolled binary wire format sized against LoRa's real
~200-byte payload budget.

Like [`toolshed`](/projects/toolshed/), this is a blueprint plus a
complete, tested vertical slice — not a finished product. No real LoRa
radio was available to test against, so the protocol engine runs against
a simulated, controllable radio-range model; a real radio driver is a
documented, swappable extension point (see
[`meshfeed/BLUEPRINT.md`](https://github.com/isidman/USPD/blob/main/meshfeed/BLUEPRINT.md)),
along with what's deliberately not built yet: message signing and
duty-cycle-aware transmission pacing.

Run `go run ./cmd/demo` from the `meshfeed/` directory to see a subscriber
get nothing while out of range, then sync automatically through an
intermediate relay node — that never itself subscribed to anything — once
in range.

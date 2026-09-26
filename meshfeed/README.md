# meshfeed

Subscribe to a node on a LoRa mesh; get its updates automatically whenever
you're in range — directly or through a relay. Read
[`BLUEPRINT.md`](BLUEPRINT.md) first: it explains the protocol, the prior
art it's built on (Secure Scuttlebutt's feed model, Meshtastic's flood
routing), and what's deliberately not built yet (signing, real radio
hardware).

No real LoRa radio is available in the environment this was built in, so
everything runs against a simulated transport that lets you control which
nodes are "in range" of each other — see `internal/simtransport`.

## Running the demo

Requires Go 1.24+, no external dependencies.

```sh
go test ./...      # run the test suite
go run ./cmd/demo    # runs the three-node scenario, then serves the result at :8080
```

The demo prints a subscriber getting nothing while out of range, then
syncing automatically through an intermediate relay node once in range —
the exact behavior this blueprint exists to prove.

## Running the reader UI

Requires Node 18+. Run `go run ./cmd/demo` first (it serves on :8080), then:

```sh
cd web
npm install
npm test        # vitest suite
npm run dev       # dev server on :5173, proxying /feeds to :8080
npm run build      # type-checks and produces web/dist
```

## API reference

Read-only — subscribing is node configuration, not a web action.

| Method | Path                   | Returns                          |
|--------|------------------------|-----------------------------------|
| GET    | `/feeds`               | `[{id, latest_seq}]`               |
| GET    | `/feeds/{id}/items`    | `[{seq, timestamp, content}]`       |

# toolshed

A working reference implementation of a generalized tool-library lending
platform — tools, hardware, and software resources, one shared
checkout/return workflow. Read [`BLUEPRINT.md`](BLUEPRINT.md) first: it
explains what this is (a tested vertical slice, not a finished product) and
how to extend it.

## Running the backend

Requires Go 1.24+. No external dependencies — nothing to download beyond
Go itself.

```sh
go test ./...          # run the test suite
go run ./cmd/server     # starts the API on :8080, storing data in ./resources.json and ./loans.json
```

Environment variables:

- `TOOLSHED_ADDR` — listen address (default `:8080`)
- `TOOLSHED_DATA_DIR` — where `resources.json` / `loans.json` are written (default `.`)

## Running the frontend

Requires Node 18+.

```sh
cd web
npm install
npm test          # runs the vitest suite
npm run dev        # dev server on :5173, proxying /resources and /loans to :8080
npm run build       # type-checks and produces web/dist
```

Run the backend (`go run ./cmd/server`, default port 8080) alongside
`npm run dev` — the frontend expects it there.

## Trying it without the UI

```sh
curl -X POST localhost:8080/resources -d '{"kind":"tool","name":"Cordless Drill"}'
curl localhost:8080/resources
curl -X POST localhost:8080/resources/<id>/checkout -d '{"borrower_id":"alice"}'
curl -X POST localhost:8080/loans/<loan-id>/return
```

See [`examples/`](examples/) for a scripted version of this flow.

## API reference

| Method | Path                          | Body                                          |
|--------|-------------------------------|------------------------------------------------|
| POST   | `/resources`                  | `{kind, name, description?, metadata?}`         |
| GET    | `/resources`                  | —                                                |
| POST   | `/resources/{id}/checkout`    | `{borrower_id, duration_hours?}` (default 168h) |
| POST   | `/loans/{id}/return`          | —                                                |

`kind` is one of `tool`, `hardware`, `software`. `GET /resources` includes
a computed `available` boolean per resource.

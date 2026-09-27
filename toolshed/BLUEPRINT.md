# toolshed: a blueprint for the "WordPress of tool libraries"

## The idea

Tool libraries (Toronto Tool Library, Ottawa Tool Library, dozens of small
community ones) keep independently building the same software: a catalog
of things a community owns collectively, a checkout/return workflow, and a
rule that says "you can't borrow what's already out." Every one of them
hand-rolls this because nothing open and reusable has become the default —
the way WordPress became the default for "I need a website with a blog."

This directory is a working answer to "what would that default look like,"
generalized one step further: the thing being shared doesn't have to be a
physical tool. It can be:

- a **tool** (a drill, a ladder, a dehydrator)
- a piece of **hardware** (a 3D printer, a diagnostic kit, a Raspberry Pi
  cluster a hackerspace lends out)
- a **software** resource (a staging server seat, a licensed account, a
  maintained fork someone's giving access to)

All three share the exact same lending mechanic: one item, one open loan at
a time, checked out and returned. That mechanic is the actual product.
Everything else — richer metadata per kind, reservations, late fees, a
public catalog page — is a layer on top of it.

## What's actually built here

This is **not** the full platform. Building a production-grade multi-tenant
app with auth, a real database, an admin UI, and a polished public catalog
in one sitting would mean either a shallow, untested pile of code, or
quietly shrinking the scope anyway — see `CLAUDE.md` at the repo root, rule
"no half-finished implementations."

What's here instead is a **complete, tested, working vertical slice** of
the one mechanic that matters, in both languages asked for:

- `internal/domain` — the `Resource` and `Loan` types. `Kind` is a string
  (`tool` / `hardware` / `software`), not three separate Go types, and
  each resource carries an open `Metadata` map for kind-specific detail.
  This is CLAUDE.md rule 3 (parameterize, don't hardcode) applied directly:
  adding a fourth kind — a seed library, a book — is a validation-list
  change, not a new code path.
- `internal/lending` — the one business rule (`CheckOut` fails if a
  resource has an open loan; `Return` closes it) against two repository
  interfaces it doesn't know the implementation of.
- `internal/storage/memory` and `internal/storage/jsonfile` — two
  interchangeable implementations of those interfaces. `memory` backs the
  tests; `jsonfile` backs the runnable demo server, storing data as
  human-readable JSON with zero database to install. **Swap in Postgres or
  SQLite by implementing the same two interfaces — `lending` and `api`
  don't change.** That's rule 4 (modular) and rule 5 (isolate what's
  likely to break) made concrete, not just asserted.
- `internal/api` — a thin HTTP layer (stdlib `net/http`, Go 1.22+'s
  built-in method+path routing, zero router dependency) translating
  requests into `lending.Service` calls.
- `web/` — a TypeScript frontend with **no framework**: plain DOM
  rendering, a typed `fetch` client, one pure formatting module. Someone
  extending this doesn't need to learn React first (rule 7: learning-curve
  budget).
- Tests at every layer: Go unit tests for the business rule (in-memory
  storage, no network), Go integration tests for the HTTP API
  (`httptest`, real requests), TypeScript unit tests for both the pure
  formatting logic and the API client (mocked `fetch`).

Every one of these was actually run — `go test ./...`, `npm run build`,
`vitest run`, and a live curl/Playwright pass against the running server —
not asserted to "probably work." See the PR this shipped in for the output.

## The zero-dependency choice

The Go module has no third-party dependencies at all — not even a SQLite
driver. This was deliberate, not an oversight: a blueprint that pulls in
five libraries teaches "here's how to integrate five libraries," not "here's
how tool lending works." `internal/storage/jsonfile` proves the storage
seam works without needing anything beyond the standard library; a real
deployment swapping in Postgres pays that dependency cost once, on purpose,
not as a side effect of following this blueprint.

## What a real deployment adds on top

None of this is built here — each is a self-contained addition behind the
same seams:

1. **A real database backend.** Implement `lending.ResourceRepository` and
   `lending.LoanRepository` against Postgres or SQLite. `jsonfile` is the
   template to follow; nothing in `lending` or `api` changes.
2. **Auth.** Wrap `api.Handler.Routes()` with middleware that populates a
   borrower ID from a session instead of trusting a client-supplied
   `borrower_id` string (today's version trusts the caller — fine for a
   single-community internal tool, not fine once it's public-facing).
3. **Reservations / waitlists.** A `Reservation` type alongside `Loan`, and
   a rule in `lending.Service` for "next in line gets notified on return."
4. **Kind-specific views.** The frontend already branches on `kind`
   (`kindLabel`); richer per-kind display (a software resource showing its
   URL, a hardware item showing its serial number from `Metadata`) is a
   frontend-only change.
5. **Multi-community / federation.** Each community's toolshed instance is
   independent by default here (like Mastodon instances). Federating
   catalogs across instances is a real design problem — worth its own
   RFC, not bolted on speculatively here (rule 6: no extraneous
   components before they're needed).

## Adding a new kind of resource

1. Add the string constant to `domain.Kind` in `internal/domain/resource.go`
   and to its `Valid()` switch.
2. Add it to the `<select>` options in `web/src/index.html` (or wherever
   the real frontend's creation form lives) and to `kindLabel` in
   `web/src/format.ts`.
3. That's it for the lending mechanic — `lending.Service` and
   `internal/api` don't know or care what kinds exist; they only see
   `domain.Kind` as an opaque validated string.

## Running it

See `README.md` in this directory.

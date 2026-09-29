---
title: "Isolate what's likely to break: the repository pattern"
summary: "Put storage, network calls, and anything else that changes out from under an interface, so swapping it never touches your business logic."
difficulty: intermediate
related_projects:
  - toolshed
---

## The problem

Storage backends change. A prototype starts with an in-memory map; a real
deployment needs a database; a test needs something fast and disposable.
If your business logic calls `sql.Query(...)` directly, every one of those
changes means editing the code that actually matters — the part that
decides *what* to store, not *where*.

USPD's own [principles]({{ '/principles/' | relative_url }}) name this
directly: isolate what's likely to break behind a clearly named, swappable
adapter. "Likely to break" isn't limited to literal failures — it also
means "likely to change," which storage always does.

## The pattern

Define the interface where you *use* it, not where you implement it. This
is the opposite of how it's often taught (define the interface next to the
struct that implements it) — but if the interface lives with the consumer,
the consumer never needs to know which implementation it's talking to.

In pseudocode:

```
// In the package that does the actual work:
interface ResourceRepository {
    Create(resource) error
    Get(id) (resource, error)
    List() ([]resource, error)
}

type Service struct {
    resources ResourceRepository   // not a concrete database type
}

func (s *Service) AddResource(r resource) error {
    // business logic here never mentions SQL, JSON files, or memory maps
    return s.resources.Create(r)
}
```

Then write as many implementations of that interface as you need:

```
// memory/store.go — for tests, zero setup
type Store struct { data map[string]resource }
func (s *Store) Create(r resource) error { s.data[r.ID] = r; return nil }

// jsonfile/store.go — for a real but simple deployment
type Store struct { path string }
func (s *Store) Create(r resource) error { /* read file, append, write file */ }
```

`Service` doesn't change no matter which one you pass it.

## Worked example

[`toolshed`]({{ '/projects/toolshed/' | relative_url }}) does exactly
this. `lending.ResourceRepository` and `lending.LoanRepository` are
defined in `internal/lending/service.go` — next to `lending.Service`,
the thing that *uses* them, not next to any implementation.

Two implementations exist, both passing the same tests because they
satisfy the same interface:

- `internal/storage/memory` — an in-memory map, used by every unit test.
  No setup, no cleanup, no network.
- `internal/storage/jsonfile` — a real but simple backend for the runnable
  demo, storing data as human-readable JSON on disk.

`toolshed/BLUEPRINT.md` names the next step this pattern sets up for free:
"swap in Postgres or SQLite by implementing the same two interfaces —
`lending` and `api` don't change." Nobody has done that yet, but the
seam is already there, proven by the fact that two implementations
already exist and neither required touching `Service`.

## How to tell if you've done this right

Try to describe your business logic out loud without saying the name of
a database, file format, or network library. If you can't, the interface
boundary is in the wrong place — it's leaking implementation detail into
the part of the code that should only care about *what*, not *how*.

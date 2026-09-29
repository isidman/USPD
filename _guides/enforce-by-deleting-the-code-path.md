---
title: "Enforce a rule by deleting the code path, not adding a check"
summary: "The most reliable way to make an operation impossible is to never write the function that performs it."
difficulty: beginner
related_projects:
  - deltos
---

## Two ways to forbid something

Say your system needs a rule like "users can't transfer credits to each
other." There are two ways to build that:

**Option A — a permission check.** Write a `transfer(from, to, amount)`
function, then guard it: `if !currentUser.canTransfer() { return error }`.
The capability exists in the code; something else decides, at runtime,
whether to allow it this time.

**Option B — no function.** Don't write `transfer` at all. There is no
code path that moves value between two accounts, anywhere in the module.

Option A is usually where people start, because it maps directly onto how
the rule is described ("users can't transfer" sounds like "check if they
can, then decide"). But every permission check is something that can be
misconfigured, bypassed by a new code path that forgets to call it, or
quietly changed by someone who doesn't know why it's there. A function
that was never written can't be bypassed, because there's nothing to call.

This is [USPD principle 6]({{ '/principles/' | relative_url }}) applied
to correctness, not just dependencies: "the most reliable component is
the one that doesn't exist." A permission check is a component. Its
absence is more reliable than its presence.

## When this applies

Not every rule can be an absence — sometimes the capability has to exist
for some users and not others, and a permission check is the only honest
option. Option B works specifically when **nobody**, under any
circumstance, is supposed to have the capability. If the rule is really
"only admins can transfer," that's Option A — a check, applied correctly.
If the rule is "credits are not a currency, they cannot be transferred,
full stop," that's Option B — and writing the function "just in case,
gated behind a check" undersells how firm the rule actually is.

## Worked example

[`deltos`]({{ '/projects/deltos/' | relative_url }})'s credit ledger has
exactly this rule: credits are non-transferable by design, stated as an
architectural constraint in the project's own specification, not a policy
decision that could reasonably change. `src/engine/ledger.ts`'s `Ledger`
class has one way in (`record`) and several ways to read (`history`,
`currentBalance`, and two projection helpers) — but critically, no way to
move value between two ledgers:

```typescript
class Ledger {
  record(input: NewContribution): Contribution { /* ... */ }
  currentBalance(now: number, essentialsOptedOut: boolean): number { /* ... */ }
  // ...history() and two projection methods, all read-only
}
```

No `transfer`, `send`, or `spend` method exists anywhere in the file, or
anywhere in the package. The test suite documents this as deliberate
rather than accidental — `ledger.test.ts` asserts those methods are
`undefined` on the class, so if one is ever added, the test fails and
points back at the comment explaining why it wasn't there in the first
place, rather than silently allowing a rule to erode.

## The tell

If you find yourself writing a comment like `// never call this in
production` or `// this should always return false`, that's a sign the
capability shouldn't exist as a callable function at all. Delete the
function; keep the comment as the explanation for why it's gone.

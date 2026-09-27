# deltos: a community contribution ledger, boiled down

## What this is

Deltos is a mobile-first app built for the Figma Config Makeathon: a
community contribution ledger that replaces market prices with real-cost
tracking. Name from δέλτος — the ancient Greek writing tablet/ledger.
Inspired by Solarpunk economics, the Hopamine Green Hackathon, and the
Integral Collective's federated cooperative economy architecture.

The original is a full mobile app (onboarding wizard, dashboard with live
charts, on-device OCR identity verification, a weighted trust score,
peer vouching, a community review queue, an 18-article help center, even
an easter-egg tapping game) built with Figma Design + Figma Make. That
source doesn't exist in this repository — this is a **from-scratch,
boiled-down reconstruction of the one mechanic that makes it interesting**:
the credit ledger and its decay, built from the project's own written
specification, not a port of the Figma Make output.

Same discipline as `toolshed/` and `meshfeed/`: a blueprint for the whole
idea, a complete and tested implementation of its core, and an honest list
of everything else that isn't built here.

## The core idea

Replace price with real cost. A market price hides what something actually
took to make — the ecological cost of a plastic toy isn't on its price
tag. Deltos makes that cost visible in the credit itself:

```
net_credits = labor_hours + (materials_kg × 0.5) − (energy_kWh × 0.3) − (co2_kg × 0.2)
```

You earn credits for labor and materials contributed; you lose credits for
the energy and CO2 spent doing it. Two people doing "the same" task don't
earn the same amount if one of them did it more wastefully — the
ecological cost is priced into the credit, not hidden in a market rate
nobody sees.

## Decay: why credits aren't money

Credits lose 10% of their value every week they sit unused (15% if the
account has opted out of Community Essentials — shared tools, communal
energy). This isn't a novel idea invented for Deltos: it's
**[demurrage currency](https://en.wikipedia.org/wiki/Demurrage_(currency))**,
the same principle behind Silvio Gesell's *Freigeld* and real regional
currencies like the Chiemgauer — money engineered to be spent or
reinvested rather than hoarded, because holding it costs you. Deltos
applies the same mechanic to contribution credits instead of cash, for the
same reason Solarpunk rejects wealth concentration: a ledger you can
"save up" indefinitely re-creates the exact accumulation dynamic a
real-cost economy is trying to escape.

The opt-out penalty isn't a separate punishment bolted onto the
economics — it's the same decay mechanism at a steeper rate, nudging
toward shared infrastructure (commons) over going it alone (private
accumulation), which is the same Solarpunk-City principle CLAUDE.md's own
rules are built from.

## Non-transferability, enforced by absence

The original spec states this as an architectural constraint: credits
can't be sent from one person to another. This implementation enforces it
the same way CLAUDE.md rule 6 recommends handling anything you want to be
impossible: **there is no transfer method**. `Ledger` in
`src/engine/ledger.ts` has exactly two operations — `record()` a new
contribution, and `currentBalance()` to read the total. There is no
`send()`, `spend()`, or `transfer()` anywhere in this package. That's not
a permission check guarding a capability that exists; the capability
was simply never written. The most reliable way to prevent an operation
is for its code path not to exist.

## Never stored, always recomputed

The original handover doc calls this out as a deliberate, debugged-into
decision: "balance always calculated live, never stored." This
implementation follows it literally — `Ledger` has no `balance` field
anywhere. `currentBalance(now, essentialsOptedOut)` walks the full
contribution history and recomputes the decayed sum every time it's
called. There's nothing to go stale, and nothing a missed decay tick could
forget to update, because there is no persisted number in the first place.

## What's actually built here

- `src/engine/formula.ts` — the net-credits calculation.
- `src/engine/decay.ts` — the exponential decay factor, at both the
  standard and Community-Essentials-opted-out rates.
- `src/engine/ledger.ts` — the ledger itself: record, read history,
  compute live balance, and two projection functions (`projectDecayOnly`
  for "what happens if I contribute nothing more," `projectWithRate` for
  "what happens if I keep contributing at roughly this rate") — the two
  traces the original app's Logs & Decay chart plotted.
- 21 unit tests across all three files, checked against hand-computed
  values (e.g. 10 credits after 2 weeks of standard decay is exactly
  10 × 0.9² = 8.1 — asserted directly, not just smoke-tested).
- A minimal, framework-free UI (`src/main.ts`) — add a contribution,
  toggle the Community Essentials opt-out, see the live balance, each
  contribution's raw vs. currently-decayed value, and both projections
  update immediately.

## Known simplification: the opt-out setting is static

The engine treats `essentialsOptedOut` as the account's *current* setting
and applies it uniformly across the entire history when computing balance.
The original design implies someone could opt in and out over time, which
would mean different weeks of history decaying at different rates — that's
not modeled here. Handling a time-varying rate would mean storing the
opt-out state *change points*, not just a boolean, and recomputing decay
piecewise across them. Worth doing before this becomes more than a
reference implementation; not done here because it roughly doubles the
complexity of `currentBalance` for a case the original handover doc never
actually specifies the behavior for.

## What's deliberately not built here

The full app's real product surface — onboarding, the dashboard UI,
on-device OCR verification, the weighted trust score, peer vouching, the
community review queue, the help center, node connectivity — none of it
is reconstructed here. Each is a real, separately-scoped feature; bolting
shallow versions of all of them onto this repository would produce exactly
the "half-finished implementation" CLAUDE.md rule 6 warns against. This
blueprint's job is the one mechanic worth getting precisely right.

**The Integral Collective alignment (CDS/OAD/COS/ITC/FRS) is explicitly
parked, not guessed at.** The original project's own handover notes that
work on reconciling Deltos's data model against the Integral Collective's
dev guide schemas was paused pending actually reading both source
documents — the right call, and one this blueprint respects rather than
overriding. Imposing an ITC-shaped schema on this engine now, without that
reading having happened, would mean guessing at a decision that isn't
this blueprint's to make. If and when that alignment work resumes, the
natural integration point is `src/engine/types.ts` — `Contribution` is
the shape that would map onto ITC's `LaborEvent` /
`MaterialConsumptionEvent`, and `Ledger` onto `ITCAccount` /
`ITCLedgerEntry`.

## Extending the formula

Adding a new cost or credit factor (say, water usage) is a two-line
change: add the field to `Contribution` in `types.ts`, and add its
weighted term to `netCredits` in `formula.ts`. Nothing in `decay.ts` or
`ledger.ts` needs to know a new factor exists — they operate on whatever
`netCredits` returns, not on the individual fields that produced it.

## Running it

See `README.md` in this directory.

# deltos

A boiled-down, tested reference implementation of the credit ledger and
decay mechanic from Deltos, a community contribution ledger built for the
Figma Config Makeathon. Read [`BLUEPRINT.md`](BLUEPRINT.md) first — it
explains the formula, why credits decay (demurrage currency, not a novel
invention), how non-transferability is enforced by a method simply not
existing, and what's deliberately not rebuilt here.

## Running it

Requires Node 18+.

```sh
npm install
npm test        # 21 unit tests across formula, decay, and ledger
npm run dev       # dev server — add contributions, toggle Community
                    # Essentials, watch the live balance and projections
npm run build      # type-checks and produces dist/
```

## The core API

```ts
import { Ledger } from "./src/engine/ledger";

const ledger = new Ledger();
ledger.record({
  description: "Repaired the shared water pump",
  category: "repairs",
  laborHours: 6,
  materialsKg: 2,
  energyKwh: 1,
  co2Kg: 0.5,
});

ledger.currentBalance(Date.now(), /* essentialsOptedOut */ false);
ledger.projectDecayOnly(Date.now(), /* weeksAhead */ 4, false);
ledger.projectWithRate(Date.now(), 4, /* assumedWeeklyRate */ 5, false);
```

There is no `transfer`, `send`, or `spend` method — see `BLUEPRINT.md` for
why that's the point, not a gap.

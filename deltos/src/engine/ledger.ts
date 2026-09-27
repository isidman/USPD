import { decayFactor, weeklyDecayRate, weeksBetween, WEEK_MS } from "./decay";
import { netCredits } from "./formula";
import type { Contribution, NewContribution } from "./types";

let nextId = 1;
function generateId(): string {
  return `contribution-${nextId++}`;
}

// Ledger holds one identity's contribution history. Read it end to end:
// there is no transfer(), send(), or spend() method here, in this file,
// or anywhere in this package. That absence *is* the non-transferability
// rule — not a permission check something could bypass, but a capability
// that was simply never written. The most reliable way to prevent an
// operation is for the code path to not exist.
//
// Known simplification: `essentialsOptedOut` is passed in at query time as
// the account's *current* setting and applied uniformly across the whole
// history. The original design implies a person could opt in and out over
// time, which would mean different weeks decaying at different rates — not
// modeled here. See BLUEPRINT.md.
export class Ledger {
  private contributions: Contribution[] = [];

  record(input: NewContribution, timestamp = Date.now()): Contribution {
    const contribution: Contribution = { ...input, id: generateId(), timestamp };
    this.contributions.push(contribution);
    return contribution;
  }

  history(): readonly Contribution[] {
    return this.contributions;
  }

  // Never stored, always recomputed: there is no `this.balance` field
  // anywhere in this class for a bug to leave stale or a decay tick to
  // forget to update.
  currentBalance(now: number, essentialsOptedOut: boolean): number {
    return this.contributions.reduce((sum, c) => {
      const weeks = weeksBetween(c.timestamp, now);
      return sum + netCredits(c) * decayFactor(weeks, essentialsOptedOut);
    }, 0);
  }

  // Balance at a future point assuming no new contributions — pure decay
  // of what's already on the books. This is the "actual" trace on the
  // original app's Logs & Decay chart, extended forward.
  projectDecayOnly(now: number, weeksAhead: number, essentialsOptedOut: boolean): number {
    return this.currentBalance(now + weeksAhead * WEEK_MS, essentialsOptedOut);
  }

  // Balance at a future point assuming a constant additional weekly net
  // credit rate on top of decay of existing history — the "projected"
  // dashed line, driven by the original UI's extra-hours-per-week slider.
  // Modeled as a contribution of `weeklyNetCreditsRate` arriving at the end
  // of each future week, each decaying independently from its own arrival
  // point to the projection horizon.
  projectWithRate(
    now: number,
    weeksAhead: number,
    weeklyNetCreditsRate: number,
    essentialsOptedOut: boolean
  ): number {
    const decayedExisting = this.projectDecayOnly(now, weeksAhead, essentialsOptedOut);
    const rate = weeklyDecayRate(essentialsOptedOut);
    let futureContribution = 0;
    for (let weekArrived = 1; weekArrived <= weeksAhead; weekArrived++) {
      const weeksToDecay = weeksAhead - weekArrived;
      futureContribution += weeklyNetCreditsRate * Math.pow(1 - rate, weeksToDecay);
    }
    return decayedExisting + futureContribution;
  }
}

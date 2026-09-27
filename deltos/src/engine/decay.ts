// Decay is what stops the ledger from becoming a bank account. Every
// credit loses value the longer it sits unused — the real-world analogue
// is demurrage currency (Silvio Gesell's Freigeld, the Chiemgauer regional
// currency): money that's designed to be spent or reinvested, not hoarded,
// because holding it costs you. Deltos applies the same idea to
// contribution credits instead of cash.

export const WEEK_MS = 7 * 24 * 60 * 60 * 1000;

export const DECAY_RATE_STANDARD = 0.1; // 10% per week
// Opting out of Community Essentials (shared tools, communal energy)
// raises the decay rate — a nudge toward participating in shared
// infrastructure rather than going it alone, not a punishment bolted on
// separately from the economics.
export const DECAY_RATE_ESSENTIALS_OPTED_OUT = 0.15;

export function weeklyDecayRate(essentialsOptedOut: boolean): number {
  return essentialsOptedOut ? DECAY_RATE_ESSENTIALS_OPTED_OUT : DECAY_RATE_STANDARD;
}

export function weeksBetween(fromMs: number, toMs: number): number {
  return Math.max(0, (toMs - fromMs) / WEEK_MS);
}

// The value remaining of one credit after `weeksElapsed` weeks of decay —
// exponential, not linear: it approaches zero but never quite reaches it,
// same shape as radioactive decay or compound interest running backward.
export function decayFactor(weeksElapsed: number, essentialsOptedOut: boolean): number {
  const rate = weeklyDecayRate(essentialsOptedOut);
  return Math.pow(1 - rate, weeksElapsed);
}

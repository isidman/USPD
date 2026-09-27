import { describe, expect, it } from "vitest";
import { WEEK_MS } from "./decay";
import { Ledger } from "./ledger";

const START = Date.parse("2026-01-01T00:00:00Z");

function contribution(overrides: Partial<Parameters<Ledger["record"]>[0]> = {}) {
  return {
    description: "test",
    category: "labor",
    laborHours: 0,
    materialsKg: 0,
    energyKwh: 0,
    co2Kg: 0,
    ...overrides,
  };
}

describe("Ledger.currentBalance", () => {
  it("is zero with no history", () => {
    const ledger = new Ledger();
    expect(ledger.currentBalance(START, false)).toBe(0);
  });

  it("equals the raw net credits the moment a contribution is recorded", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 5 }), START);
    expect(ledger.currentBalance(START, false)).toBeCloseTo(5);
  });

  it("decays a single contribution week over week", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 10 }), START);

    expect(ledger.currentBalance(START + WEEK_MS, false)).toBeCloseTo(9);
    expect(ledger.currentBalance(START + 2 * WEEK_MS, false)).toBeCloseTo(8.1);
  });

    it("sums multiple contributions, each decaying independently from its own timestamp", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 10 }), START); // 2 weeks old at `now`
    ledger.record(contribution({ laborHours: 10 }), START + WEEK_MS); // 1 week old at `now`

    const now = START + 2 * WEEK_MS;
    // first: 10 * 0.9^2 = 8.1; second: 10 * 0.9^1 = 9
    expect(ledger.currentBalance(now, false)).toBeCloseTo(8.1 + 9);
  });

  it("decays faster for an account opted out of Community Essentials", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 10 }), START);
    const now = START + WEEK_MS;

    const optedIn = ledger.currentBalance(now, false);
    const optedOut = ledger.currentBalance(now, true);
    expect(optedOut).toBeLessThan(optedIn);
  });

  it("has no method that moves credits between ledgers", () => {
    const ledger = new Ledger();
    // If a transfer/send/spend method is ever added, this test documents
    // that its absence was deliberate — see the class-level comment in
    // ledger.ts. Failing this test is a signal to re-read that comment
    // before adding one.
    expect((ledger as unknown as Record<string, unknown>).transfer).toBeUndefined();
    expect((ledger as unknown as Record<string, unknown>).send).toBeUndefined();
    expect((ledger as unknown as Record<string, unknown>).spend).toBeUndefined();
  });
});

describe("Ledger.projectDecayOnly", () => {
  it("matches currentBalance computed directly at the future timestamp", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 10 }), START);

    const viaProjection = ledger.projectDecayOnly(START, 3, false);
    const viaDirect = ledger.currentBalance(START + 3 * WEEK_MS, false);
    expect(viaProjection).toBeCloseTo(viaDirect);
  });
});

describe("Ledger.projectWithRate", () => {
  it("matches projectDecayOnly when the assumed future rate is zero", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 10 }), START);

    expect(ledger.projectWithRate(START, 4, 0, false)).toBeCloseTo(ledger.projectDecayOnly(START, 4, false));
  });

  it("adds decayed future contributions on top of decayed existing balance", () => {
    const ledger = new Ledger();
    // One week ahead: one future contribution arrives and is evaluated
    // immediately at the horizon, so it hasn't decayed yet.
    const result = ledger.projectWithRate(START, 1, 5, false);
    expect(result).toBeCloseTo(5);
  });

  it("a higher assumed weekly rate always projects to a higher balance", () => {
    const ledger = new Ledger();
    ledger.record(contribution({ laborHours: 10 }), START);

    const low = ledger.projectWithRate(START, 8, 1, false);
    const high = ledger.projectWithRate(START, 8, 5, false);
    expect(high).toBeGreaterThan(low);
  });
});

import { describe, expect, it } from "vitest";
import { decayFactor, weeksBetween, WEEK_MS } from "./decay";

describe("decayFactor", () => {
  it("is 1 at zero weeks elapsed — nothing decayed yet", () => {
    expect(decayFactor(0, false)).toBe(1);
  });

  it("applies the standard 10% weekly rate", () => {
    expect(decayFactor(1, false)).toBeCloseTo(0.9);
    expect(decayFactor(2, false)).toBeCloseTo(0.81);
  });

  it("applies the steeper 15% rate when Community Essentials is opted out", () => {
    expect(decayFactor(1, true)).toBeCloseTo(0.85);
    expect(decayFactor(2, true)).toBeCloseTo(0.7225);
  });

  it("decays faster when opted out than when opted in, at the same elapsed time", () => {
    expect(decayFactor(4, true)).toBeLessThan(decayFactor(4, false));
  });

  it("approaches but never reaches zero", () => {
    expect(decayFactor(100, false)).toBeGreaterThan(0);
  });
});

describe("weeksBetween", () => {
  it("computes fractional weeks from millisecond timestamps", () => {
    const start = 0;
    const end = WEEK_MS * 2.5;
    expect(weeksBetween(start, end)).toBeCloseTo(2.5);
  });

  it("never returns negative — a future timestamp clamps to zero elapsed", () => {
    expect(weeksBetween(WEEK_MS, 0)).toBe(0);
  });
});

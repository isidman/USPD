import { describe, expect, it } from "vitest";
import { netCredits } from "./formula";

describe("netCredits", () => {
  it("credits labor and materials, at full and half weight respectively", () => {
    const result = netCredits({ laborHours: 4, materialsKg: 2, energyKwh: 0, co2Kg: 0 });
    expect(result).toBeCloseTo(4 + 2 * 0.5);
  });

  it("deducts energy and CO2 costs", () => {
    const result = netCredits({ laborHours: 0, materialsKg: 0, energyKwh: 10, co2Kg: 5 });
    expect(result).toBeCloseTo(-(10 * 0.3) - 5 * 0.2);
  });

  it("can go negative when ecological cost exceeds labor and materials given", () => {
    const result = netCredits({ laborHours: 1, materialsKg: 0, energyKwh: 20, co2Kg: 10 });
    expect(result).toBeLessThan(0);
  });

  it("matches the documented example exactly", () => {
    // Net = labor_hours + (materials_kg × 0.5) - (energy_kWh × 0.3) - (co2_kg × 0.2)
    const result = netCredits({ laborHours: 3, materialsKg: 4, energyKwh: 2, co2Kg: 1 });
    expect(result).toBeCloseTo(3 + 4 * 0.5 - 2 * 0.3 - 1 * 0.2);
  });
});

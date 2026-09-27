import type { Contribution } from "./types";

// The weights that make ecological cost visible instead of hidden in a
// market price. A contribution earns credits for labor and materials
// given, and loses credits for energy and CO2 spent producing it — so two
// people doing "the same" task don't earn the same amount if one of them
// did it more wastefully.
export const WEIGHTS = {
  materials: 0.5,
  energy: 0.3,
  co2: 0.2,
} as const;

type CostFields = Pick<Contribution, "laborHours" | "materialsKg" | "energyKwh" | "co2Kg">;

export function netCredits(c: CostFields): number {
  return c.laborHours + c.materialsKg * WEIGHTS.materials - c.energyKwh * WEIGHTS.energy - c.co2Kg * WEIGHTS.co2;
}

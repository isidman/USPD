// A Contribution is the only fact this engine ever stores. Everything else
// — balance, decay, projections — is derived from a list of these, every
// time it's asked for. There's no separate "account balance" record
// anywhere to drift out of sync with the history that justifies it.
export interface Contribution {
  id: string;
  timestamp: number; // unix ms
  description: string;
  category: string;
  laborHours: number;
  materialsKg: number;
  energyKwh: number;
  co2Kg: number;
}

export type NewContribution = Omit<Contribution, "id" | "timestamp">;

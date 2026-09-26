// Mirrors toolshed/internal/domain — kept as plain types (no framework,
// no runtime validation library) so anyone reading this as a blueprint
// doesn't need to learn a validation library just to see the shape of the
// data (CLAUDE.md rule 7: learning-curve budget).

export type Kind = "tool" | "hardware" | "software";

export interface Resource {
  id: string;
  kind: Kind;
  name: string;
  description: string;
  metadata?: Record<string, string>;
  created_at: string;
}

export interface ResourceView extends Resource {
  available: boolean;
}

export interface Loan {
  id: string;
  resource_id: string;
  borrower_id: string;
  checked_out_at: string;
  due_at: string;
  returned_at?: string;
}

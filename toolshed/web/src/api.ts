import type { Kind, Loan, ResourceView } from "./types";

// Every function here does exactly one HTTP call and returns typed data —
// this is the seam a real app would mock in tests, and the only file that
// knows the API's base URL or request shapes.

const BASE_URL = (import.meta as { env?: { VITE_API_BASE?: string } }).env?.VITE_API_BASE ?? "";

async function asJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(body.error ?? `request failed with status ${response.status}`);
  }
  return response.json() as Promise<T>;
}

export async function listResources(): Promise<ResourceView[]> {
  const res = await fetch(`${BASE_URL}/resources`);
  return asJSON<ResourceView[]>(res);
}

export interface CreateResourceInput {
  kind: Kind;
  name: string;
  description?: string;
  metadata?: Record<string, string>;
}

export async function createResource(input: CreateResourceInput): Promise<ResourceView> {
  const res = await fetch(`${BASE_URL}/resources`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  const created = await asJSON<ResourceView>(res);
  return { ...created, available: true };
}

export async function checkOut(resourceId: string, borrowerId: string, durationHours?: number): Promise<Loan> {
  const res = await fetch(`${BASE_URL}/resources/${resourceId}/checkout`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ borrower_id: borrowerId, duration_hours: durationHours }),
  });
  return asJSON<Loan>(res);
}

export async function returnLoan(loanId: string): Promise<Loan> {
  const res = await fetch(`${BASE_URL}/loans/${loanId}/return`, { method: "POST" });
  return asJSON<Loan>(res);
}

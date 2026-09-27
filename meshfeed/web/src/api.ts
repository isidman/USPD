import type { FeedView, ItemView } from "./types";

const BASE_URL = (import.meta as { env?: { VITE_API_BASE?: string } }).env?.VITE_API_BASE ?? "";

async function asJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new Error(`request failed with status ${response.status}`);
  }
  return response.json() as Promise<T>;
}

export async function listFeeds(): Promise<FeedView[]> {
  const res = await fetch(`${BASE_URL}/feeds`);
  return asJSON<FeedView[]>(res);
}

export async function listItems(feedId: string): Promise<ItemView[]> {
  const res = await fetch(`${BASE_URL}/feeds/${feedId}/items`);
  return asJSON<ItemView[]>(res);
}

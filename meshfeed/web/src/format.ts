// Pure formatting, kept separate from api.ts (network) and main.ts (DOM)
// so it's trivial to test — same split toolshed's frontend uses.

export function formatTimestamp(unixSeconds: number): string {
  if (unixSeconds === 0) return "unknown time";
  return new Date(unixSeconds * 1000).toISOString();
}

export function shortFeedId(id: string): string {
  return id.length > 12 ? `${id.slice(0, 12)}…` : id;
}

import { afterEach, describe, expect, it, vi } from "vitest";
import { listFeeds, listItems } from "./api";

function jsonResponse(body: unknown, ok = true, status = 200): Response {
  return { ok, status, json: async () => body } as Response;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("listFeeds", () => {
  it("calls GET /feeds", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse([{ id: "abc", latest_seq: 3 }]));
    vi.stubGlobal("fetch", fetchMock);

    const feeds = await listFeeds();

    expect(fetchMock).toHaveBeenCalledWith("/feeds");
    expect(feeds).toEqual([{ id: "abc", latest_seq: 3 }]);
  });

  it("throws on a non-ok response", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse(null, false, 500)));
    await expect(listFeeds()).rejects.toThrow("500");
  });
});

describe("listItems", () => {
  it("calls GET /feeds/:id/items", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse([{ seq: 1, timestamp: 100, content: "hi" }]));
    vi.stubGlobal("fetch", fetchMock);

    const items = await listItems("abc");

    expect(fetchMock).toHaveBeenCalledWith("/feeds/abc/items");
    expect(items[0].content).toBe("hi");
  });
});

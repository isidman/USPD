import { afterEach, describe, expect, it, vi } from "vitest";
import { checkOut, listResources } from "./api";

function jsonResponse(body: unknown, ok = true, status = 200): Response {
  return {
    ok,
    status,
    json: async () => body,
  } as Response;
}

describe("listResources", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("calls GET /resources and returns parsed JSON", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse([{ id: "1", kind: "tool", name: "Drill", description: "", created_at: "now", available: true }])
    );
    vi.stubGlobal("fetch", fetchMock);

    const resources = await listResources();

    expect(fetchMock).toHaveBeenCalledWith("/resources");
    expect(resources).toHaveLength(1);
    expect(resources[0].name).toBe("Drill");
  });

  it("throws with the server's error message on a non-ok response", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse({ error: "boom" }, false, 500)));

    await expect(listResources()).rejects.toThrow("boom");
  });
});

describe("checkOut", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("posts borrower_id and duration_hours to the checkout endpoint", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse({ id: "loan-1", resource_id: "r1", borrower_id: "alice", checked_out_at: "now", due_at: "later" })
    );
    vi.stubGlobal("fetch", fetchMock);

    const loan = await checkOut("r1", "alice", 48);

    expect(fetchMock).toHaveBeenCalledWith(
      "/resources/r1/checkout",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ borrower_id: "alice", duration_hours: 48 }),
      })
    );
    expect(loan.borrower_id).toBe("alice");
  });
});

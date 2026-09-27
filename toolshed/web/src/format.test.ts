import { describe, expect, it } from "vitest";
import { describeAvailability, kindLabel } from "./format";
import type { ResourceView } from "./types";

function baseResource(overrides: Partial<ResourceView> = {}): ResourceView {
  return {
    id: "abc",
    kind: "tool",
    name: "Drill",
    description: "",
    created_at: "2026-01-01T00:00:00Z",
    available: true,
    ...overrides,
  };
}

describe("describeAvailability", () => {
  it("reports available resources as available", () => {
    expect(describeAvailability(baseResource({ available: true }))).toBe("available");
  });

  it("reports checked-out resources as checked out", () => {
    expect(describeAvailability(baseResource({ available: false }))).toBe("checked out");
  });
});

describe("kindLabel", () => {
  it("labels each kind", () => {
    expect(kindLabel("tool")).toBe("Tool");
    expect(kindLabel("hardware")).toBe("Hardware");
    expect(kindLabel("software")).toBe("Software");
  });
});

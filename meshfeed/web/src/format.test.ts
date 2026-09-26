import { describe, expect, it } from "vitest";
import { formatTimestamp, shortFeedId } from "./format";

describe("formatTimestamp", () => {
  it("formats a unix timestamp as ISO", () => {
    expect(formatTimestamp(1_700_000_000)).toBe(new Date(1_700_000_000_000).toISOString());
  });

  it("reports unknown for a zero timestamp", () => {
    expect(formatTimestamp(0)).toBe("unknown time");
  });
});

describe("shortFeedId", () => {
  it("truncates long ids", () => {
    expect(shortFeedId("357a59c7a22abde9")).toBe("357a59c7a22a…");
  });

  it("leaves short ids alone", () => {
    expect(shortFeedId("abc")).toBe("abc");
  });
});

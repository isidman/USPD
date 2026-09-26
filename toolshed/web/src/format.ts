import type { ResourceView } from "./types";

// Pure functions are the easiest thing in a codebase to test and the
// easiest to trust — kept separate from api.ts (which touches the
// network) and main.ts (which touches the DOM) on purpose.

export function describeAvailability(resource: ResourceView): string {
  return resource.available ? "available" : "checked out";
}

export function kindLabel(kind: ResourceView["kind"]): string {
  switch (kind) {
    case "tool":
      return "Tool";
    case "hardware":
      return "Hardware";
    case "software":
      return "Software";
  }
}

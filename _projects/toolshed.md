---
title: "toolshed"
tagline: "A blueprint and working reference implementation for a generalized tool-lending platform: tools, hardware, and software, one shared checkout workflow."
category: tools
status: idea
license: "AGPL-3.0-or-later (code), CC BY-SA 4.0 (docs)"
repo_url: "https://github.com/isidman/USPD/tree/main/toolshed"
docs_url: "https://github.com/isidman/USPD/blob/main/toolshed/BLUEPRINT.md"
tech_stack:
  - "Go (backend, zero third-party dependencies)"
  - "TypeScript (frontend, no framework)"
principles:
  repairable: true
  common_components: true
  parameterized: true
  modular: true
  isolates_failure: true
  minimal: true
  learnable: true
  legible: true
---

Tool libraries (Toronto, Ottawa, dozens of small community ones) keep
independently rebuilding the same software: a catalog of shared things, a
checkout/return workflow, and a rule against double-booking. Nothing open
has become the default the way WordPress became the default for websites.

toolshed generalizes the idea one step further — the thing being shared
doesn't have to be a physical tool. It can be a **tool**, a piece of
**hardware** (a 3D printer, a diagnostic kit), or a **software** resource
(a staging server seat, a maintained fork someone's giving access to). All
three share one lending mechanic: one item, one open loan at a time.

This isn't a finished platform — it's a complete, tested vertical slice of
that one mechanic, in both Go and TypeScript, plus a blueprint document
explaining how to extend it into a real deployment (a real database, auth,
reservations, kind-specific views). See
[`toolshed/BLUEPRINT.md`](https://github.com/isidman/USPD/blob/main/toolshed/BLUEPRINT.md)
for the full writeup, and [`toolshed/README.md`](https://github.com/isidman/USPD/blob/main/toolshed/README.md)
to run it locally — it's zero-dependency on the Go side, so there's nothing
to install beyond Go itself.

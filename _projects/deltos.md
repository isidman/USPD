---
title: "deltos"
tagline: "A community contribution ledger that replaces market price with real cost — labor and materials earn credits, energy and CO2 spend them, and unused balance decays like demurrage currency."
category: community
status: idea
license: "AGPL-3.0-or-later (code), CC BY-SA 4.0 (docs)"
repo_url: "https://github.com/isidman/USPD/tree/main/deltos"
docs_url: "https://github.com/isidman/USPD/blob/main/deltos/BLUEPRINT.md"
tech_stack:
  - "TypeScript (pure calculation engine, zero dependencies, no backend)"
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

Deltos was built for the Figma Config Makeathon: a mobile-first community
contribution ledger inspired by Solarpunk economics, the Hopamine Green
Hackathon, and the Integral Collective's cooperative economy architecture.
The full app (onboarding, on-device OCR verification, a weighted trust
score, peer vouching, a help center) was built with Figma Design + Figma
Make — that source isn't in this repository.

This entry is a from-scratch, boiled-down reconstruction of the one
mechanic worth getting precisely right: `net_credits = labor_hours +
(materials_kg × 0.5) − (energy_kWh × 0.3) − (co2_kg × 0.2)`, with weekly
exponential decay preventing accumulation — the same principle behind
real demurrage currencies like Gesell's Freigeld — and non-transferability
enforced the way [`toolshed`](/projects/toolshed/) and
[`meshfeed`](/projects/meshfeed/) enforce their own constraints: by the
capability simply not existing in the code, not by a permission check.

See [`deltos/BLUEPRINT.md`](https://github.com/isidman/USPD/blob/main/deltos/BLUEPRINT.md)
for the full writeup, including what's deliberately not rebuilt (the
app's UI surface) and why the project's own paused work on aligning to
the Integral Collective's ITC schema is respected here rather than
guessed at.

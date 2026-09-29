---
title: Principles
permalink: /principles/
---

# Why we build this way

USPD catalogs and documents software built along Solarpunk lines: technology
that serves people and ecosystems, designed to be repaired and remixed
rather than replaced. The 8 rules below are USPD's own answer to *how* to
build that way — for what Solarpunk itself actually is, as a movement,
see [Further reading](#further-reading) at the bottom of this page.

They're also the checklist we ask contributors to run a project against
before cataloging it here (see the `principles` field in a
[project entry]({{ '/contributing/' | relative_url }})), and the same
checklist this repository's own code is written against — see
[`CLAUDE.md`](https://github.com/isidman/USPD/blob/main/CLAUDE.md) in the
repo root.

1. **Repair over replace.** Prefer mature, widely-maintained dependencies.
   A broken or abandoned dependency should be swappable, not a reason to
   rewrite the project.

2. **Common and bulk over novel.** Default to boring, mainstream tools.
   Exotic tech should earn its place, not be the default.

3. **Parameterize, don't hardcode.** Configuration adapts the software to
   a new context; the code itself shouldn't need to change.

4. **Modular, not monolithic.** Clear boundaries between components mean
   any one part can be replaced without touching the rest.

5. **Isolate what's likely to break.** Calls to the outside world (network,
   third-party APIs) live behind a clearly named, swappable adapter.

6. **No extraneous components.** No dependency for one convenience
   function, no abstraction for a hypothetical future.

7. **Respect the learning-curve budget.** A contributor should have to
   learn at most one new high-value thing to get involved.

8. **Optimize for the reader's mental model.** Obvious naming, one obvious
   way to do a thing — a newcomer should understand a module in minutes.

## Community

- Docs explain what a project is and who it's for, in plain language.
- Contribution is welcome by default; the bar for a first contribution
  stays low.
- Projects listed here should stay forkable — permissive or copyleft
  licensing, never closed.

## Sources

The 8 numbered rules above come from three places, cross-checked against
each other rather than taken from just one:

- **Pete Seeger's repair hierarchy** — reduce, reuse, repair, rebuild,
  refurbish, refinish, resell, recycle, compost, in that order of
  preference. Rules 1 and 6 are this applied to dependencies: prefer
  keeping what exists over replacing it, and prefer not needing the thing
  at all over adding it.
- **[Solarpunk City](https://solarpunkcity.org/)**'s design principles
  (CC0 — no rights reserved). Rules 2, 3, 4, and 7 (common/bulk
  components, parameterized designs, modularity, and designs a builder
  can actually implement) are taken close to verbatim from its published
  guidance for physical and hardware projects, adapted here for software.
  Its Makerspace-practice writeup is the direct source for the Community
  section above.
- **The [GOSH Manifesto](https://openhardware.science/gosh-manifesto/)**
  (Gathering for Open Science Hardware) — independently arrived at
  criteria for open hardware (repairable, lowest-cost materials,
  understandable documentation, non-territorial collaboration) that
  overlap heavily with Solarpunk City's, which is why both are cited: two
  unrelated communities converging on the same rules is a stronger signal
  than either alone.

Rule 6 ("no extraneous components") and the general framing of software
sustainability also draw on current **green/sustainable software
engineering** work — the [Green Software Foundation](https://greensoftware.foundation/)'s
practices, the GREENS workshop series on energy-efficient software design,
and the LIMITS workshop on computing degrowth. This repository's own
[`CLAUDE.md`](https://github.com/isidman/USPD/blob/main/CLAUDE.md) names
these explicitly as the source cross-checked against Seeger/Solarpunk
City/GOSH when the rules were first written down.

**Individual catalog entries cite further prior art of their own** where a
specific design decision has a direct real-world precedent, rather than
repeating it here at the site level:

- [`meshfeed`]({{ '/projects/meshfeed/' | relative_url }})'s sync protocol
  is explicitly built on [Secure Scuttlebutt](https://en.wikipedia.org/wiki/Secure_Scuttlebutt)'s
  append-only feed/gossip model and [Meshtastic](https://meshtastic.org/docs/overview/mesh-algo/)'s
  managed flood-routing mechanism.
- [`deltos`]({{ '/projects/deltos/' | relative_url }})'s decaying credit
  balance is explicitly framed as a software implementation of
  **demurrage currency** (Silvio Gesell's *Freigeld*, and real regional
  currencies like the Chiemgauer) rather than a novel invention.

See each project's own `BLUEPRINT.md` for the full reasoning.

## Further reading

Everything above is USPD's own engineering answer to *how* to build
Solarpunk software. It isn't where the word "Solarpunk" itself comes
from. If you want the movement's own account of what it is and why,
these predate and outrank anything on this page:

- **[Solarpunk: Notes toward a Manifesto](https://hieroglyph.asu.edu/2014/09/solarpunk-notes-toward-a-manifesto/)**
  (Adam Flynn, 2014) — the essay generally credited with coining and
  crystallizing the term.
- **[Solarpunk: A Reference Guide](https://medium.com/solarpunks/solarpunk-a-reference-guide-8bcf18871965)**
  (Jay Springett, 2017, published on [solarpunks.net](https://solarpunks.net/))
  — a curated compendium across architecture, technology, and economics
  (circular and gift economies, alternative currencies, decentralized
  production) that shaped how the movement talks about itself. Its
  economics thread is the same lineage `deltos`'s credit system sits in.
- **[A Solarpunk Manifesto](https://re-des.org/a-solarpunk-manifesto/)**
  (2019) — a later synthesis explicitly built from the two essays above,
  and the version most often cited today.

These three are about Solarpunk as a movement — fiction, aesthetics,
politics, hope. This site is only about one narrow slice of that: how to
write software that doesn't betray it.

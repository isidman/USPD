---
title: Principles
permalink: /principles/
---

# Why we build this way

USPD catalogs and documents software built along Solarpunk lines: technology
that serves people and ecosystems, designed to be repaired and remixed
rather than replaced. These principles come from Pete Seeger's repair
hierarchy, [Solarpunk City](https://solarpunkcity.org/)'s design principles
(CC0), and the [GOSH Manifesto](https://openhardware.science/gosh-manifesto/)
for open hardware — adapted for software.

They're also the checklist we ask contributors to run a project against
before cataloging it here (see the `principles` field in a
[project entry]({{ '/contributing/' | relative_url }})).

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

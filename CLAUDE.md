# USPD — Engineering Principles

This project follows Solarpunk / GOSH-style design principles. Check every
change against this list before writing it, not after.

Source: Pete Seeger's repair hierarchy + Solarpunk City (solarpunkcity.org,
CC0) design principles, cross-checked against the GOSH (Gathering for Open
Source Hardware) Manifesto and current sustainable/green software engineering
practice (Green Software Foundation, GREENS workshop, LIMITS workshop on
computing degrowth).

## Rules

1. **Repair over replace.** Prefer stdlib and mature, widely-maintained
   dependencies over exotic ones. If a dependency disappears or breaks,
   swapping it must be a local, contained change — never a rewrite.

2. **Common and bulk over novel.** Default to boring, mainstream tech
   (language, framework, format). Don't introduce a niche tool unless it's
   the one high-value thing worth the learning cost (see rule 7).

3. **Parameterize, don't hardcode.** Behavior comes from config, not magic
   numbers or baked-in paths/URLs. A new deployment context should need
   config changes, not code changes.

4. **Modular, not monolithic.** Components have explicit interfaces and
   clear boundaries. Any single piece can be cut out and replaced without
   touching the rest.

5. **Isolate what's likely to break.** Anything touching the outside world
   (network calls, third-party APIs, external services) sits behind one
   clearly named adapter — easy to find, easy to swap, easy to mock in
   tests.

6. **No extraneous components.** No dependency for one convenience
   function. No abstraction for a hypothetical future. The most reliable
   component is the one that doesn't exist. (Matches: don't add
   speculative flexibility, don't build past what's asked.)

7. **Respect the learning-curve budget.** If a contributor has to learn
   something new to work on this, it should be one high-value thing
   (e.g. "learn Docker"), not a stack of bespoke tools and DSLs.

8. **Optimize for the reader's mental model, not completeness.** Obvious
   naming, one obvious way to do a thing. A newcomer should be able to
   build a correct-enough mental model of any module in minutes.

## Community conventions (from Makerspace practice)

- README and docs explain what this is and who it's for in plain terms —
  no generic boilerplate. If someone asks a question the README should
  have answered, that's a signal to fix the README.
- Contribution is welcome by default; keep barriers to a first contribution
  low.
- License stays permissive/forkable — if governance or maintainership ever
  goes bad, the ability to fork cleanly is the community's exit valve.

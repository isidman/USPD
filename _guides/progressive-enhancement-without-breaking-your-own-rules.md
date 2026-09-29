---
title: "Add a dependency without breaking your own zero-dependency rule"
summary: "When a real request conflicts with one of your own principles, keep the principle's intent alive through how you build the exception, not by ignoring the request."
difficulty: intermediate
---

## The situation

USPD's [principles]({{ '/principles/' | relative_url }}) say "no
extraneous components" and "common and bulk over novel." For a while, that
meant this site shipped zero JavaScript — plain HTML and CSS, nothing to
load, nothing to fail. Then a real request came in: add animations, using
a specific library, for a better feel.

That's a direct request, not a suggestion to weigh — but "just add the
script tag" would quietly break the actual reason the zero-JS rule
existed: content that depends on a third-party CDN to become visible is
one ad-blocker or one outage away from a blank page. The fix isn't to
refuse the request or to follow it carelessly. It's to build the
exception so the *reason* behind the old rule still holds, even though
the rule's letter no longer does.

## The method: progressive enhancement, actually verified

"Progressive enhancement" is a familiar phrase — build a working baseline,
then layer enhancements on top that degrade safely if they fail. The part
usually skipped is the *verified* half: most code that claims to degrade
gracefully has never actually been run with the enhancement missing.

1. **Make the baseline state the real, final state.** Not a hidden state
   waiting for JavaScript to reveal it — the plain-CSS version should be
   what a user is looking at if nothing else ever loads.
2. **Let the enhancement layer add on top, never replace.** It should be
   physically impossible for the enhancement's failure to leave the page
   in a worse state than the baseline.
3. **Test the failure case, not just the success case.** Block the
   dependency. Set the accessibility preference that should disable it.
   Confirm the baseline still holds under both — with an actual browser,
   not by inspecting the code and assuming it's fine.

## Worked example

USPD's own homepage does this for its card entrance animation. The
CSS (`assets/css/style.css`) never hides `.reveal` elements — they're
fully visible by default, full stop. `assets/js/site.js` is the entire
enhancement, and it's structured as a sequence of early-exit guards:

```js
if (typeof anime === "undefined") return;               // CDN failed to load
if (matchMedia("(prefers-reduced-motion: reduce)").matches) return; // visitor opted out
// only past both guards does anything get animated —
// and even then, elements were already visible before this ran
```

This was verified three separate ways with a real browser (Playwright),
not asserted in a comment and trusted:

- The CDN request blocked entirely → content stayed fully visible.
- `prefers-reduced-motion: reduce` set → content stayed fully visible
  immediately, no animation delay.
- The real library actually loaded → the animation genuinely ran, caught
  mid-interpolation (opacity at 0.89 partway through a fade) before
  settling at 1.

All three had to pass before the change shipped. Two of three are the
"boring" failure cases — the ones that are easy to skip testing because
nothing looks broken in the editor.

## The general lesson

A principle being overridden by a specific request isn't a reason to drop
the principle everywhere it touches the change. It's a reason to ask what
the principle was actually protecting against, and build the exception so
that protection survives — then prove it survived, rather than assuming
good structure implies correct behavior.

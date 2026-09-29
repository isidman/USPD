// Progressive enhancement only. Every element this touches is already
// fully visible via plain CSS — this script's job is to make things that
// are already there ease in, never to reveal things that were hidden
// waiting for JS. If anime.js fails to load (CDN down, ad blocker, no
// JS at all), the page looks and works exactly the same, just without
// the entrance motion.
(function () {
  if (typeof anime === "undefined") return;

  var reducedMotion = window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  if (reducedMotion) return;

  var targets = document.querySelectorAll(".reveal");
  if (targets.length === 0) return;

  // Set the "from" state and animate to "to" back-to-back, both after the
  // real content has already rendered normally — there is no window in
  // which the content is invisible waiting on this script.
  anime.set(targets, { opacity: 0, translateY: 14 });
  anime({
    targets: targets,
    opacity: [0, 1],
    translateY: [14, 0],
    delay: anime.stagger(70),
    duration: 520,
    easing: "easeOutQuad",
  });
})();

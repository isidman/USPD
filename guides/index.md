---
title: Guides
permalink: /guides/
---

# Guides

How-to writeups for building solarpunk-aligned software: patterns, tutorials,
and lessons learned.

{% assign guides = site.guides | sort: "title" %}
{% for guide in guides %}
<div class="card reveal">
  <h3><a href="{{ guide.url | relative_url }}">{{ guide.title }}</a></h3>
  <p class="meta">{% if guide.difficulty %}<span class="badge" data-difficulty="{{ guide.difficulty }}">{{ guide.difficulty }}</span>{% endif %}</p>
  <p>{{ guide.summary }}</p>
</div>
{% else %}
<p>No guides written yet. <a href="/contributing/">Write the first one.</a></p>
{% endfor %}

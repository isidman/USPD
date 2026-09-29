---
title: Home
---

# Projects

USPD catalogs open source projects built along solarpunk lines: repairable,
modular, common-components-first, and legible to newcomers. See
[Principles]({{ '/principles/' | relative_url }}) for what that means in practice.

{% assign projects = site.projects | sort: "title" %}
{% for project in projects %}
<div class="card reveal">
  <h3><a href="{{ project.url | relative_url }}">{{ project.title }}</a></h3>
  <p class="meta">Tags: {{ project.status }}, {{ project.category }}</p>
  <p>{{ project.tagline }}</p>
</div>
{% else %}
<p>No projects catalogued yet. <a href="{{ '/contributing/' | relative_url }}">Add the first one.</a></p>
{% endfor %}

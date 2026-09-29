---
title: Contribute
permalink: /contributing/
---

# Contributing to USPD

USPD grows through pull requests. Two kinds of entries: **projects**
(software that exists, cataloged here) and **guides** (how-to writeups).

## Add a project

1. Copy [`templates/project-entry.md`](https://github.com/isidman/USPD/blob/main/templates/project-entry.md)
   into `_projects/<your-project-slug>.md`.
2. Fill in the frontmatter. Be honest on the `principles` checklist — a
   project with a few `false` entries is still welcome; the checklist is
   for readers deciding whether to build on it, not a gate.
3. Write a short body: what it does, why it exists, how to get started.
4. Open a pull request.

## Add a guide

1. Copy [`templates/guide-entry.md`](https://github.com/isidman/USPD/blob/main/templates/guide-entry.md)
   into `_guides/<your-guide-slug>.md`.
2. Fill in the frontmatter and write the guide.
3. Open a pull request.

## Running the site locally

This is a plain Jekyll site — no build step beyond Jekyll itself.

```
bundle install
bundle exec jekyll serve
```

Then open `http://localhost:4000`.

## Licensing

By contributing, you agree your writeups (project entries, guides) are
licensed under [CC BY-SA 4.0]({{ '/LICENSE-CONTENT' | relative_url }}), and any code you
contribute to the site itself (layouts, config, scripts) is licensed under
[AGPL-3.0-or-later]({{ '/LICENSE-CODE' | relative_url }}) — matching the rest of the repository.

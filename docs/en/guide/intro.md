# Introduction

This project is a bilingual documentation template built on **Docsify**. Chinese and English content are maintained independently under separate directories, distinguished by URL hash prefixes (`/zh-CN/` and `/en/`).

## How It Works

```
User clicks [中文] button
    ↓
location.hash changes to #/zh-CN/...
    ↓
Docsify loads .md files from /zh-CN/ directory via alias rules
    ↓
Sidebar is also loaded from /zh-CN/_sidebar.md
```

## Why Directory Separation?

Compared to toggling languages with HTML comments inside the same Markdown file, **directory separation** offers:

- ✅ Each language maintained independently, no interference
- ✅ Sidebar structure can be customized per language
- ✅ Accurate search results, no cross-language pollution
- ✅ Easy for teams to collaborate on translations

[→ Quick Start](/en/guide/quickstart)

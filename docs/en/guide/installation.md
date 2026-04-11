# Installation

## Requirements

- Node.js ≥ 14 (only if using docsify-cli)
- Any modern browser

## Install docsify-cli globally

```bash
npm install -g docsify-cli
docsify --version
```

## Initialize a new project

```bash
docsify init ./docs
```

This generates in `./docs`:

```
docs/
├── index.html
├── README.md
└── .nojekyll    ← Required for GitHub Pages to serve underscore files
```

## Use without CLI

Just include the CDN script in your HTML — no `npm install` or build step needed:

```html
<script src="//cdn.jsdelivr.net/npm/docsify@4/lib/docsify.min.js"></script>
```

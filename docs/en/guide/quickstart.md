# Quick Start

## 1. Clone the repo

```bash
git clone https://github.com/your/repo.git
cd repo
```

## 2. Start a local server

Docsify requires an HTTP server (you can't just open the HTML file directly).

**Option A: docsify-cli (recommended)**

```bash
npm install -g docsify-cli
docsify serve docs
# Visit http://localhost:3000
```

**Option B: Python**

```bash
python -m http.server 3000
# Visit http://localhost:3000
```

**Option C: VS Code Live Server**

Install the [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) extension, right-click `index.html` → Open with Live Server.

## 3. Add new content

Create a `.md` file in the appropriate language directory and add a link in `_sidebar.md`:

```
en/
└── guide/
    └── my-new-page.md   ← new file
```

```markdown
<!-- en/_sidebar.md -->
* [My New Page](/en/guide/my-new-page)
```

## 4. Deploy

Upload the entire directory to any static hosting platform:

- GitHub Pages
- Vercel / Netlify
- AWS S3 / Cloudflare Pages

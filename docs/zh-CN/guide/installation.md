# 安装

## 系统要求

- Node.js ≥ 14（如需使用 docsify-cli）
- 任意现代浏览器

## 全局安装 docsify-cli

```bash
npm install -g docsify-cli
docsify --version
```

## 初始化新项目

```bash
docsify init ./docs
```

这会在 `./docs` 目录下生成：

```
docs/
├── index.html
├── README.md
└── .nojekyll    ← GitHub Pages 必需，防止忽略下划线文件
```

## 在现有项目中引入

如果不使用 CLI，只需在 HTML 中引入 CDN 即可：

```html
<script src="//cdn.jsdelivr.net/npm/docsify@4/lib/docsify.min.js"></script>
```

无需 `npm install`，无需构建工具。

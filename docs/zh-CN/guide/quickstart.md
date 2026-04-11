# 快速开始

## 1. 克隆项目

```bash
git clone https://github.com/your/repo.git
cd repo
```

## 2. 启动本地服务器

Docsify 需要一个 HTTP 服务器（不能直接双击 HTML 打开）。

**方式一：使用 docsify-cli（推荐）**

```bash
npm install -g docsify-cli
docsify serve docs
# 访问 http://localhost:3000
```

**方式二：使用 Python**

```bash
# Python 3
python -m http.server 3000
# 访问 http://localhost:3000
```

**方式三：使用 VS Code 插件**

安装 [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) 插件，右键 `index.html` → Open with Live Server。

## 3. 新增内容

在对应语言目录下创建 `.md` 文件，再在 `_sidebar.md` 中添加链接即可：

```
zh-CN/
└── guide/
    └── my-new-page.md   ← 新建文件
```

```markdown
<!-- zh-CN/_sidebar.md -->
* [我的新页面](/zh-CN/guide/my-new-page)
```

## 4. 部署

将整个目录上传到任意静态托管平台：

- GitHub Pages
- Vercel / Netlify
- 阿里云 OSS / 腾讯云 COS

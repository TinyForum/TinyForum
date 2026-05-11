package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PluginBrowser(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>插件目录浏览器 - TinyForum</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 1000px; margin: 0 auto; padding: 20px; background: #f5f5f5; }
        .container { background: white; border-radius: 8px; padding: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        h1 { font-size: 1.5rem; margin-top: 0; }
        .breadcrumb { margin-bottom: 20px; font-size: 0.9rem; word-break: break-all; }
        .breadcrumb a { color: #3498db; text-decoration: none; }
        .breadcrumb a:hover { text-decoration: underline; }
        ul { list-style: none; padding-left: 0; }
        li { margin: 8px 0; padding: 4px 0; border-bottom: 1px solid #eee; }
        .dir, .file { cursor: pointer; display: inline-block; }
        .dir { color: #2c3e50; font-weight: 500; }
        .file { color: #2980b9; }
        .dir:hover, .file:hover { text-decoration: underline; }
        .size { color: #888; margin-left: 10px; font-size: 0.8rem; }
        .loading { color: #666; font-style: italic; }
        .error { color: #e74c3c; }
        footer { margin-top: 30px; text-align: center; font-size: 0.8rem; color: #777; }
        .modal { display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; overflow: auto; background-color: rgba(0,0,0,0.5); }
        .modal-content { background-color: #fefefe; margin: 5% auto; padding: 20px; border-radius: 8px; width: 80%; max-width: 900px; box-shadow: 0 4px 20px rgba(0,0,0,0.2); }
        .close { color: #aaa; float: right; font-size: 28px; font-weight: bold; cursor: pointer; }
        .close:hover { color: black; }
        .markdown-body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 16px; line-height: 1.5; word-wrap: break-word; }
        .markdown-body h1, .markdown-body h2 { border-bottom: 1px solid #eaecef; }
        .markdown-body code { background: #f6f8fa; padding: 0.2em 0.4em; border-radius: 3px; }
        .markdown-body pre { background: #f6f8fa; padding: 16px; overflow: auto; border-radius: 3px; }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.1.0/github-markdown.min.css">
</head>
<body>
<div class="container">
    <h1>📁 插件目录浏览器</h1>
    <div class="breadcrumb" id="breadcrumb">加载中...</div>
    <div id="file-list">加载中...</div>
    <footer>TinyForum · 安全目录浏览 · <a href="/api/v1/health">健康检查</a></footer>
</div>

<div id="fileModal" class="modal">
    <div class="modal-content">
        <span class="close">&times;</span>
        <div id="modal-body" class="markdown-body">加载中...</div>
    </div>
</div>

<script>
    let currentPath = "";
    const apiBase = "/api/v1/plugins/files";
    const fileListDiv = document.getElementById("file-list");
    const breadcrumbDiv = document.getElementById("breadcrumb");
    const modal = document.getElementById("fileModal");
    const modalBody = document.getElementById("modal-body");
    const closeSpan = document.getElementsByClassName("close")[0];

    closeSpan.onclick = function() { modal.style.display = "none"; };
    window.onclick = function(event) { if (event.target == modal) modal.style.display = "none"; };

    async function loadDirectory(path = "") {
        fileListDiv.innerHTML = '<div class="loading">加载中...</div>';
        try {
            const url = apiBase + "?path=" + encodeURIComponent(path);
            const resp = await fetch(url);
            if (!resp.ok) throw new Error("HTTP " + resp.status);
            const data = await resp.json();
            renderFiles(data.files, data.path);
            renderBreadcrumb(data.path);
            currentPath = data.path;
        } catch (err) {
            fileListDiv.innerHTML = '<div class="error">加载失败: ' + err.message + '</div>';
        }
    }

    function renderFiles(files, currentPath) {
        if (!files.length) {
            fileListDiv.innerHTML = '<div>📂 当前目录为空</div>';
            return;
        }
        const ul = document.createElement("ul");
        files.forEach(file => {
            const li = document.createElement("li");
            const icon = file.isDir ? "📁" : "📄";
            const nameSpan = document.createElement("span");
            nameSpan.className = file.isDir ? "dir" : "file";
            nameSpan.textContent = icon + " " + file.name;
            
            if (file.isDir) {
                nameSpan.style.cursor = "pointer";
                nameSpan.onclick = () => {
                    const newPath = currentPath === "." ? file.name : currentPath + "/" + file.name;
                    loadDirectory(newPath);
                };
            } else {
                const fileUrl = "/store/plugins/" + (currentPath === "." ? "" : currentPath + "/") + file.name;
                const ext = file.ext.toLowerCase();
                
                if (ext === ".md") {
                    nameSpan.style.cursor = "pointer";
                    nameSpan.onclick = () => previewMarkdown(fileUrl, file.name);
                    const badge = document.createElement("span");
                    badge.textContent = " [预览]";
                    badge.style.fontSize = "0.7rem";
                    badge.style.color = "#27ae60";
                    badge.style.marginLeft = "8px";
                    nameSpan.appendChild(badge);
                } else {
                    nameSpan.style.cursor = "pointer";
                    nameSpan.onclick = () => window.open(fileUrl, "_blank");
                }
            }
            
            if (!file.isDir && file.size) {
                const sizeSpan = document.createElement("span");
                sizeSpan.className = "size";
                sizeSpan.textContent = formatSize(file.size);
                nameSpan.appendChild(sizeSpan);
            }
            li.appendChild(nameSpan);
            ul.appendChild(li);
        });
        fileListDiv.innerHTML = "";
        fileListDiv.appendChild(ul);
    }

    function renderBreadcrumb(path) {
        if (path === "." || path === "") {
            breadcrumbDiv.innerHTML = '<a href="#" onclick="loadDirectory(\'\'); return false;">根目录</a>';
            return;
        }
        const parts = path.split("/");
        let html = '<a href="#" onclick="loadDirectory(\'\'); return false;">根目录</a>';
        let accumulated = "";
        for (let i = 0; i < parts.length; i++) {
            const part = parts[i];
            if (part === "") continue;
            accumulated += (accumulated ? "/" : "") + part;
            html += ' / <a href="#" onclick="loadDirectory(\'' + accumulated + '\'); return false;">' + part + '</a>';
        }
        breadcrumbDiv.innerHTML = html;
    }

    function formatSize(bytes) {
        if (bytes < 1024) return bytes + " B";
        if (bytes < 1048576) return (bytes / 1024).toFixed(1) + " KB";
        return (bytes / 1048576).toFixed(1) + " MB";
    }

    async function previewMarkdown(fileUrl, fileName) {
        modalBody.innerHTML = '<div class="loading">加载 Markdown 内容...</div>';
        modal.style.display = "block";
        try {
            const resp = await fetch(fileUrl);
            if (!resp.ok) throw new Error("下载失败");
            const text = await resp.text();
            const htmlContent = marked.parse(text);
            // 修复点：使用字符串拼接而不是模板字符串
            modalBody.innerHTML = '<h3>' + escapeHtml(fileName) + '</h3><hr/>' + htmlContent;
        } catch (err) {
            modalBody.innerHTML = '<div class="error">预览失败: ' + err.message + '</div>';
        }
    }

    function escapeHtml(str) {
        return str.replace(/[&<>]/g, function(m) {
            if (m === '&') return '&amp;';
            if (m === '<') return '&lt;';
            if (m === '>') return '&gt;';
            return m;
        });
    }

    loadDirectory("");
</script>
</body>
</html>
`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

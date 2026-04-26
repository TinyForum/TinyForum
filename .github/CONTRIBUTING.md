# 贡献指南

感谢你愿意为 TinyForum 贡献力量！请仔细阅读以下指南，以确保协作顺畅。

## 如何提交代码

1. **Fork 本仓库** 到你的 GitHub 账户。
2. **创建特性分支**（从 `main` 分支派生）：
   ```bash
   git checkout -b feature/your-feature-name
   ```
   分支命名建议：`feature/`, `fix/`, `docs/`, `refactor/` 等前缀。
3. **提交更改**：
   ```bash
   git commit -m "简短描述本次更改"
   ```
   提交信息请遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范（如 `feat: 添加新API`，`fix: 修复评论分页错误`）。
4. **推送到你的远程分支**：
   ```bash
   git push origin feature/your-feature-name
   ```
5. **开启 Pull Request (PR)**：
   - 前往原始仓库，点击 “New Pull Request”。
   - 清晰描述 PR 的目的、改动范围以及相关 Issue（如有）。
   - 等待 CI 检查通过及维护者审核。

## 设计原则

所有新代码和 API 端点必须遵循以下核心原则：

### SOLID 原则
- **单一职责**：每个类/模块只负责一个功能领域。
- **开闭原则**：对扩展开放，对修改封闭。优先使用继承、接口或组合，避免修改稳定代码。
- **里氏替换**：子类必须能替换父类且不破坏程序。
- **接口隔离**：接口应小而专一，避免“胖接口”。
- **依赖倒置**：依赖抽象（接口/协议），而非具体实现。

### RESTful 规范
- 使用标准 HTTP 方法：
  - `GET`：获取资源（幂等、安全）。
  - `POST`：创建新资源。
  - `PUT` / `PATCH`：完整或部分更新资源。
  - `DELETE`：删除资源。
- 资源命名使用名词复数（如 `/users`, `/posts`）。
- 使用 HTTP 状态码准确表达结果（200 OK, 201 Created, 400 Bad Request, 404 Not Found, 500 Internal Server Error 等）。
- API 响应应包含一致的 JSON 格式，必要时提供错误详情字段。

## 代码风格

本项目遵循 **[Google Python Style Guide](https://google.github.io/styleguide/pyguide.html)**。关键要点：

- **缩进**：4 个空格，禁止使用 Tab。
- **行长度**：最多 **80 字符**（文档字符串和注释可放宽至 100）。
- **命名约定**：
  - 模块名：`lowercase_with_underscores.py`
  - 类名：`CamelCase`
  - 函数/变量名：`lowercase_with_underscores`
  - 常量：`UPPER_CASE_WITH_UNDERSCORES`
- **导入**：按标准库 → 第三方库 → 本地模块分组，每组之间空一行。
- **类型注解**：所有公共函数必须包含类型注解（使用 `typing` 模块）。
- **文档字符串**：使用 Google 风格的 docstring。

### 自动化格式与检查

- **格式化**：使用 [Black](https://github.com/psf/black)（行长度设为 80）。
- **代码检查**：使用 [Flake8](https://flake8.pycqa.org/) 或 [Pylint](https://pylint.org/)（配置遵循 Google 规则）。
- **类型检查**：推荐使用 `mypy`。
- **测试**：所有新功能必须包含单元测试。运行以下命令确保全部通过：
  ```bash
  pytest
  ```
  测试覆盖率应不低于 80%。

## Pull Request 检查清单

提交 PR 前，请确认：

- [ ] 代码遵循 Google 风格且通过 Black 格式化。
- [ ] 所有测试通过（`pytest`）。
- [ ] 添加了适当的类型注解。
- [ ] 如果引入了新 API，遵循 RESTful 规范并更新了 API 文档（或 OpenAPI 规范）。
- [ ] 设计上满足 SOLID 原则（无明显的职责过重或紧耦合）。
- [ ] Commit 信息符合 Conventional Commits 格式。
- [ ] PR 描述中关联了相关 Issue（如有）。

## Issue 规范

提交 Issue 时，请遵循以下模板：

### Bug 报告
- **标题**：清晰简短，例如 `[BUG] 评论接口返回 500`
- **复现步骤**：
  1. 访问 `POST /api/comments`
  2. 请求体：`{"post_id": 1, "content": ""}`
  3. 观察结果
- **期望行为**：返回 400 错误，提示内容不能为空。
- **实际行为**：服务器 500 错误。
- **环境**：操作系统、Python 版本、浏览器（如果是前端问题）。

### 功能请求
- **标题**：`[FEATURE] 希望增加帖子点赞功能`
- **描述**：清晰说明该功能解决什么需求。
- **设计建议**（可选）：如何满足 RESTful 和 SOLID 的初步思路。

## 社区与行为准则

参与本项目即表示你同意遵守 [行为准则](CODE_OF_CONDUCT.md)。请保持友善、专业，尊重他人。

---

感谢你的贡献！🎉
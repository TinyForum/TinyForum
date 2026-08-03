# 脚本工具

TinyForum 使用 `dev-script/` 目录管理所有构建、测试、部署脚本。

---

## Makefile 入口

根目录执行 `make help` 查看所有可用目标：

```bash
make help
```

### 常用命令

```bash
make init-dev        # 初始化开发环境（检查依赖、生成配置）
make run-backend     # 启动后端开发服务器
make run-frontend    # 启动前端开发服务器
make build           # 编译后端二进制
make test            # 运行所有测试
make check           # 代码质量检查（lint + format + type-check）
make code-gen        # 生成代码（Wire 依赖注入 + Swagger 文档）
make clean           # 清理构建产物
```

### 后端专用

```bash
cd backend
make run             # 启动服务
make build           # 编译
make wire            # 重新生成依赖注入代码
make docs-api        # 生成 Swagger 文档
make tidy            # 整理 go.mod 依赖
```

### 前端专用

```bash
cd frontend
pnpm dev             # 开发模式
pnpm build           # 生产构建
pnpm type-check      # TypeScript 类型检查
pnpm lint            # ESLint 检查
pnpm format          # Prettier 格式化
```

---

## 脚本目录结构

```
dev-script/
├── scripts/
│   ├── dev/          # 开发环境脚本
│   ├── env/          # 环境变量解析
│   ├── db/           # 数据库操作
│   └── nginx/        # Nginx 配置
├── Makefile.common   # 公共变量
├── Makefile.dev      # 开发环境
├── Makefile.check    # 代码检查
├── Makefile.code     # 代码生成
└── Makefile.docker   # Docker 构建
```

---

## 环境变量

项目使用 `.env` 文件管理环境变量：

```bash
# 复制示例配置
cp .env.example .env

# 编辑 .env 文件填入实际值
```

后端配置文件位于 `backend/config/`，支持 `SIGHUP` 热更新。

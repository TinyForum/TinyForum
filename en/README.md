# TinyForum Documentation

[中文文档](/zh-CN/README)

TinyForum is an open-source modern forum system built with **Go (Gin + GORM) backend** and **Next.js frontend**.

## Quick Links

| Section | Description |
|---------|-------------|
| [Getting Started](/en/guide/intro) | Installation and deployment |
| [Development Guide](/en/dev/architecture) | Architecture and coding standards |
| [API Documentation](/en/dev/swagger) | REST API reference |

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.24+, Gin, GORM, PostgreSQL, Redis |
| Frontend | Next.js 16, React 19, TypeScript, TailwindCSS |
| Auth | JWT + Casbin RBAC |
| Infrastructure | Docker Compose, Nginx |

## Features

- Multi-type content: posts, articles, Q&A, topic discussions
- Rich text editor (Tiptap)
- Nested comments with likes
- User follow & timeline
- Role-based access control (RBAC)
- Content moderation with AI-assisted sensitive word detection
- Bot automation (Lua scripts + no-code flow engine)
- Plugin system for frontend extensibility
- i18n support (Chinese & English)

## License

MIT

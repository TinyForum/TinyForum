# Tiny Forum

> Go (Gin + GORM) 后端 × Next.js 15.5 (App Router) 前端 × PostgreSQL 持久化 × Nginx 接入

<div align="center">
  <img src="./docs/assets/logo.svg" alt="Tiny Forum Logo" width="200" />
</div>

<summary>
Tiny Forum 是一个基于 Go 和 Next.js 的开源论坛项目，旨在提供一个简单、易用的社区平台。项目采用现代技术栈，包括 Gin、GORM、Next.js、PostgreSQL 等，为用户提供丰富的功能，包括帖子发布、评论、点赞、标签、用户、风控、公告、后台管理等。
<details>
正在优化，图片暂未更新
</details>
</summary>

> 特别提示，如果需要 mock 数据库，请使用 dbmock 进行数据生成: https://github.com/caoyang2002/dbmock
> 命令：`./dbmock generate -c mock_config.yml --count 10000` 为每个表生成10000个数据，如果出现了错误，请单独生成错误的表，例如：`./dbmock generate --rows posts=1000`

## 更改说明

当前本项目已经完成基础设施的设计与构建，即将开始完善。数据库模型可能会更改，请勿存放重要数据：

如果发生数据库相关错误，请重建：

```bash
psql -U postgres -c "DROP DATABASE IF EXISTS tiny_forum; CREATE DATABASE tiny_forum;"
```

## 快速开始

> 项目已支持 docker 和 podman 安装。这里仅演示 docker 安装，更多安装方式请查看文档，或者使用 `make help` 命令

> 在 Mac 上，需要先启动初始化虚拟机 `make podman-init`，然后执行 `make podman-build`。目前仍然存在一些问题，正在处理。

# 开发概述

该项目配置了 vscode 工作区，建议使用 Vscode 打开，并根据提示，启动工作区（如果没有提示，可以打开文件 `.vscode/TinyForum.code-workspace` 然后根据提示启动）

## 技术栈

| 层     | 技术                                                                           |
| ------ | ------------------------------------------------------------------------------ |
| 后端   | Go 1.26, Gin, GORM, Wire (手动注入), JWT, Zap                                  |
| 前端   | Next.js 16, TypeScript, Tailwind CSS, DaisyUI, TanStack Query, Zustand, Tiptap |
| 数据库 | PostgreSQL 16, Redis                                                           |
| 接入   | Nginx（容器部署）                                                              |
| 部署   | Docker + Docker Compose                                                        |

## 功能列表

- ✅ 用户注册 / 登录 / JWT 鉴权
- ✅ 发帖（帖子 / 文章 / 话题）、富文本编辑器
- ✅ 评论 & 嵌套回复
- ✅ 点赞 / 取消点赞
- ✅ 标签系统
- ✅ 关注 / 取消关注
- ✅ 积分系统 & 排行榜
- ✅ 站内消息通知
- ✅ 个人主页 / 编辑资料
- ✅ 管理后台（用户管理、封禁、置顶）
- ✅ 全文搜索（标题 & 内容）
- ✅ 风控（内容合规、行为风控）

# 上游依赖

1. 敏感词： https://github.com/konsheng/Sensitive-lexicon

# 路线图

## 全局规划

市面上有众多论坛基础设施，例如 `Discuz!`、`Discourse`、`BBS-GO`、`StarPro`、`Hadsky`、`Tecmz`、`巡云轻论坛`、`南生论坛` 等，这些论坛都是本项目的前辈，本项目也无意与其竞争，更具体来说，本项目希望从“可维护 + 易运营 + 用户共创”入手，构建一个持续的、健康的社区生态。

本项目会涉及开源版本（一个）和商业版本（五个），作为论坛/社区的核心功能，将持续开源，但是新增的商业化相关功能以及特殊优化，将采用闭源策略，以维系项目发展。

## 当前阶段

目前已经在和商业化相关的策略人员进行交流。

## 参与本项目

如果您想参与到项目的开发开源版本或者商业版本，欢迎添加我的微信：thomelgo

![wechat](./docs/assets/wechat.jpg)

本项目目前纯用爱发电，如若加入，还望三思。

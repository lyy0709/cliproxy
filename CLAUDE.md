# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

cli-proxy-go 是一个 Go 语言实现的 AI 代理服务，支持多平台（Claude、OpenAI、Gemini）账户管理和统一 API 接口。

## Development Commands

### Backend (Go)

```bash
make run              # 运行后端服务
make build            # 构建到 bin/server
make dev              # 热重载开发模式（需安装 air）
make test             # 运行测试
make fmt              # 格式化代码
make lint             # 代码检查（需安装 golangci-lint）
make deps             # 整理依赖
```

### Frontend (Vue 3)

```bash
cd web && npm install  # 安装依赖
cd web && npm run dev  # 开发模式（localhost:3000）
cd web && npm run build  # 构建
make web-build         # 构建前端并同步到后端
```

### Full Stack

```bash
make all              # 完整构建（前端 + 后端）
docker-compose up -d  # Docker 部署（MySQL + 应用）
```

## Architecture

```
cmd/server/           # 后端入口
internal/
├── handler/          # HTTP 处理器
├── middleware/       # 中间件（JWT、API Key 认证、CORS）
├── service/          # 业务逻辑层
├── repository/       # 数据访问层
├── model/            # 数据模型
├── proxy/
│   ├── adapter/      # 各平台适配器（Claude/OpenAI/Gemini/Azure/Bedrock）
│   └── scheduler/    # 调度器（账户池轮询、负载均衡）
├── cache/            # 缓存层
└── config/           # 配置管理
pkg/                  # 公共工具包（logger/response/utils）
web/                  # 前端（Vue 3 + Vite + Pinia）
configs/config.yaml   # 主配置文件
```

## Key Technical Details

**后端**: Go 1.21+, Gin, GORM (MySQL), JWT + API Key 双认证
**前端**: Vue 3 (Composition API + `<script setup>`), Vite 5, Pinia, Alova

**代理接口路由**:
- Claude: `/claude/v1/messages`
- OpenAI: `/openai/v1/chat/completions`
- Gemini: `/gemini/v1/chat`
- Responses API: `/responses` 或 `/v1/responses`

**管理接口**: `/api/admin/*`（需 JWT 认证）

## Coding Conventions

- Go: gofmt 格式化，文件名 snake_case，导出类型 CamelCase
- 前端: Vue 3 `<script setup>`，组件名 PascalCase，2 空格缩进
- 保留中文注释风格，错误信息包含上下文

## Default Credentials

- 管理员: `admin` / `admin123`
- 生产环境务必修改 `JWT_SECRET`、数据库密码

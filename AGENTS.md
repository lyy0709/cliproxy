# Repository Guidelines

## 项目结构与模块组织
- `cmd/server/` 为后端入口；`internal/` 按 handler/middleware/service/repository/model/proxy 分层。
- `pkg/` 放通用工具（日志、响应、JWT）；`configs/config.yaml` 为主配置。
- `web/` 为前端（Vue 3 + Vite）；`scripts/` 含压测脚本；`docker-compose.yml` 启动全栈。

## 构建、测试与开发命令
- 后端常用：`make run`（本地启动）、`make build`（输出 `bin/server`）、`make dev`（air 热重载）。
- 质量检查：`make fmt`（go fmt）、`make lint`（golangci-lint，需本地安装）。
- 前端常用：`cd web && npm install`、`npm run dev`、`npm run build`；`make web-build` 会把 `web/dist` 同步到 `internal/handler/dist` 并 gzip。
- Docker：`docker-compose up -d` 启动 MySQL + 服务。

## 编码风格与命名规范
- Go 使用 gofmt（Tab 缩进）；文件名多为 snake_case；导出类型/函数用 CamelCase。
- 重要逻辑请保留中文注释风格；错误信息要包含上下文。
- 前端使用 Vue 3 `<script setup>`；组件名 PascalCase（如 `Accounts.vue`），JS 采用 2 空格缩进。

## 测试指南
- 后端测试入口：`make test`（即 `go test ./...`）。
- 新增逻辑建议提供 `_test.go` 的表驱动测试；当前仓库未配置前端测试框架。

## 提交与 PR 规范
- Git 历史提交较少，未形成固定规范；建议使用简短动词开头的提交摘要（中英文均可）。
- PR 请描述变更动机、影响范围与验证方式；涉及 UI 变更请附截图。

## 配置与安全提示
- 默认账号/密码与 JWT 密钥在 `configs/config.yaml` 与 `docker-compose.yml` 中；生产环境务必修改 `JWT_SECRET`、`DB_*`、`ADMIN_*`。

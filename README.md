# cli-proxy-go

一个基于 Go 的多平台 AI 代理服务，提供统一 API 接口与 Web 管理后台。

## 功能特性

- 多平台适配：Claude / OpenAI / Azure OpenAI / Gemini / Bedrock
- 统一代理接口与平台路由
- OpenAI Responses API 兼容（/responses、/v1/responses）
- 账户池管理、负载均衡与故障转移
- API Key 权限控制与客户端过滤
- 使用统计、请求日志与操作审计
- 模型映射、错误规则与默认配置初始化

## 生产部署（Docker Compose）

### 1) 环境准备

- Docker + Docker Compose

### 2) 准备 .env

在项目根目录创建 `.env`：

```bash
MYSQL_ROOT_PASSWORD=请设置强密码
DB_NAME=cli-proxy
DB_USER=cli-proxy
DB_PASSWORD=请设置强密码
JWT_SECRET=请设置强密钥
DATA_ENCRYPTION_KEY=请设置强密钥
APP_PORT=8080
TZ=Asia/Shanghai
```

- `JWT_SECRET` 与 `DATA_ENCRYPTION_KEY` 建议使用 32 字节以上随机字符串
- `APP_PORT` / `TZ` 为可选项

### 3) （可选）调整配置文件

`configs/config.yaml` 为主配置文件，可调整：
- 服务端口 `server.port`
- 日志级别与目录 `log.*`
- 数据库连接池 `mysql.*`
- 缓存与会话策略 `cache.*`

### 4) 启动服务

```bash
docker compose up -d
```

常用运维命令：

```bash
docker compose logs -f
docker compose down
docker compose down -v
```

### 5) 访问地址

- 管理后台：`http://localhost:${APP_PORT}`
- 健康检查：`http://localhost:${APP_PORT}/health`
- API Base：`http://localhost:${APP_PORT}`

### 6) 数据与日志

- MySQL 数据：`./mysql-data`
- 日志目录：`./logs`
- 配置目录：`./configs`

## 默认管理员账号

- 用户名：`admin`
- 密码：`admin123`

首次登录后请立刻修改密码。

## 配置与环境变量

配置加载顺序：`configs/config.yaml` → 环境变量覆盖（优先级更高）。

应用支持的环境变量：

| 变量名 | 说明 |
| --- | --- |
| `JWT_SECRET` | JWT 密钥（生产必配） |
| `DATA_ENCRYPTION_KEY` | 敏感数据加密密钥（生产必配） |
| `DB_HOST` | MySQL 主机地址 |
| `DB_USER` | MySQL 用户名 |
| `DB_PASSWORD` | MySQL 密码 |
| `DB_NAME` | MySQL 数据库名 |
| `STATIC_DIR` | 自定义静态资源目录（默认 `web/dist`） |

Docker Compose 相关变量：

| 变量名 | 说明 |
| --- | --- |
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 |
| `APP_PORT` | 宿主机映射端口（默认 8080） |
| `TZ` | 时区（默认 Asia/Shanghai） |

## API 使用示例

### Claude

```bash
curl http://localhost:${APP_PORT}/claude/v1/messages \
  -H "x-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-sonnet-4-20250514",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 1024
  }'
```

### OpenAI Responses API

```bash
curl http://localhost:${APP_PORT}/responses \
  -H "x-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-5.1-codex-max",
    "input": "Write a hello world function",
    "stream": true
  }'
```

## 本地开发

环境要求：
- Go 1.24+
- MySQL 8.0+
- Node.js 20+

后端运行：

```bash
go build -o cli-proxy ./cmd/server
./cli-proxy
```

前端开发：

```bash
cd web
npm install
npm run dev
```

若需要后端直接托管前端静态资源：

```bash
cd web
npm run build
```

## 项目结构

```
cli-proxy/
├── cmd/server/          # 程序入口
├── internal/
│   ├── handler/         # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   └── proxy/           # 代理核心
│       ├── adapter/     # 各平台适配器
│       └── scheduler/   # 调度器
├── pkg/                 # 公共工具包
└── web/                 # 前端 (Vue 3 + Vite)
```

## License

MIT

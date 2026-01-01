# 前端构建阶段
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端源码
COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build

# 后端构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git make

# 复制 go mod 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 从前端构建阶段复制构建产物
COPY --from=frontend-builder /app/web/dist ./web/dist/

# 构建后端
WORKDIR /app
RUN make build

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 复制二进制文件和前端资源
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/web/dist /app/web/dist/
COPY --from=builder /app/configs /app/configs/

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./server"]

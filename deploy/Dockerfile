# ---------- 前端构建 ----------
FROM node:20-alpine AS fe-builder
WORKDIR /app

# 先复制依赖文件
COPY frontend/package*.json ./
RUN npm ci

# 再复制前端源码
COPY frontend/ .
RUN npm run build   # 输出到 /app/dist

# ---------- 后端构建 ----------
FROM golang:1.23-alpine AS be-builder
WORKDIR /app

# 先复制 go.mod/go.sum，加快依赖安装
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 再复制后端源码
COPY backend/ .

COPY --from=fe-builder /app/dist ./dist

# 编译
RUN go build -o auth-server .

# ---------- 运行阶段 ----------
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache tzdata ca-certificates

COPY --from=be-builder /app/auth-server .
COPY --from=be-builder /app/config-prod.yaml .

EXPOSE 8080
CMD ["./auth-server", "--config", "config-prod.yaml"]

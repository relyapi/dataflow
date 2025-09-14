# 第一阶段：构建 Vue 静态资源
FROM node:20-alpine AS build

# 设置工作目录
WORKDIR /app

# 安装 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# 复制依赖文件
COPY package.json pnpm-lock.yaml* ./

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制代码并构建
COPY . .
RUN pnpm run build

# 第二阶段：Nginx 提供静态文件
FROM caddy:2.10.2-alpine

# 复制 Caddy 配置文件
COPY deploy/Caddyfile /etc/caddy/Caddyfile

# 复制构建产物到 caddy 静态目录
COPY --from=build /app/dist /usr/share/caddy

EXPOSE 80
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
# 构建前端
FROM node:16-alpine AS frontend-builder
WORKDIR /app/frontend
# 复制前端项目的 package.json 和 package-lock.json
COPY frontend/package*.json ./
# 安装前端依赖
RUN npm install
# 复制前端项目的所有文件
COPY frontend ./
# 构建前端项目
RUN npm run build

# 构建后端
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/backend
# 复制后端项目的 go.mod 和 go.sum
COPY backend/go.* ./
# 下载后端依赖
RUN go mod download
# 复制后端项目的所有文件
COPY backend ./
# 编译后端项目
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 最终镜像
FROM nginx:alpine
# 设置工作目录
WORKDIR /app
# 复制前端构建的静态文件到 Nginx 的默认目录
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html
# 复制后端可执行文件到容器中
COPY --from=backend-builder /app/backend/main /app/
# 复制 Nginx 配置文件
COPY deploy/nginx/nginx.conf /etc/nginx/conf.d/default.conf

# 设置环境变量
ENV GIN_MODE=release

# 添加执行权限
RUN chmod +x /app/main

# 暴露端口
EXPOSE 80 8080

# 启动后端服务和 Nginx
CMD ["/bin/sh", "-c", "/app/main & nginx -g 'daemon off;'"]
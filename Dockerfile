# 构建前端
FROM node:16 as frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

# 构建后端
FROM golang:1.21 as backend-builder
WORKDIR /app/backend
COPY backend/go.* ./
RUN go mod download
COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 最终镜像
FROM nginx:alpine
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html
COPY --from=backend-builder /app/backend/main /usr/local/bin/
COPY deploy/nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["sh", "-c", "/usr/local/bin/main & nginx -g 'daemon off;'"]

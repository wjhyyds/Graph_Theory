# 图论算法可视化项目

本项目使用 Golang 作为后端，Vue.js 作为前端，实现了以下算法的可视化：
- Prim 算法（最小生成树）
- Dijkstra 算法（最短路径）
- Floyd 算法（多源最短路径）
- Bayes 算法（概率图模型）
- 最大匹配算法
- 旅行商问题（TSP）

## 项目结构
```
.
├── backend/         # Golang 后端代码
│   ├── algorithms/  # 算法实现
│   ├── api/        # API 处理
│   └── main.go     # 主程序
├── frontend/       # Vue.js 前端代码
│   ├── src/        # 源代码
│   └── public/     # 静态资源
└── deploy/         # 部署相关配置
    ├── nginx/      # Nginx 配置
    └── docker/     # Docker 配置
```

## 运行项目
1. 启动后端服务
```bash
cd backend
go run main.go
```

2. 启动前端服务
```bash
cd frontend
npm install
npm run serve
```

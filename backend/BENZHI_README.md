# team-task-board-service

## 标准命令

```bash
# 先在仓库根目录启动 MySQL：docker compose up -d db
cd backend
go build ./...        # 编译
go run ./cmd/server   # 启动 HTTP 服务（默认 :8080）
go test ./...         # 测试（SQLite 测试需启用 CGO）
```

健康检查：`GET /health` 或 `GET /healthz`。

## 环境

- 基础镜像: golang:1.22
- 依赖已在镜像构建阶段预下载，容器内离线可用。
- 代码目录: /app
- 运行时需要可访问的 MySQL；默认配置由环境变量 `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME` 提供。

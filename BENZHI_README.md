# team-task-board-service

## 本机验证

```bash
cp .env.example .env
docker compose up --build db backend
curl http://127.0.0.1:19503/healthz
```

也可以在已启动 MySQL 并正确设置 `DB_*` 环境变量后直接运行：

```bash
cd backend
go build ./...
go run ./cmd/server
```

## 评测镜像

```bash
docker build -f benzhi.Dockerfile -t team-task-board-service .
```

镜像会真实启动 `./cmd/server`；服务运行时需要可访问的 MySQL。完整本地依赖和健康检查请使用根目录 Compose。

# 团队任务看板

团队任务看板是一个支持拖拽操作的协作式项目管理工具，采用看板（Kanban）视图展示任务状态流转，适用于敏捷开发和小团队协作。

## 项目主要功能

- 工作区与看板管理：创建多个工作区，每个工作区包含多个看板，看板默认提供「待办 / 进行中 / 已完成」三列，支持自定义添加列。
- 拖拽任务卡片：基于 `@hello-pangea/dnd` 实现同列排序与跨列状态流转。
- 任务详情与属性：标题、Markdown 描述、负责人、彩色标签、优先级、截止日期、子任务清单、附件上传、评论记录。
- 成员协作与权限：邀请成员并设置管理员 / 编辑者 / 观察者角色；管理员管理工作区与成员，编辑者创建、编辑、移动任务，观察者只读。
- 任务筛选与搜索：按负责人、标签、优先级、截止日期筛选，关键词搜索标题与描述，匹配卡片高亮、其余卡片半透明。
- 活动日志与通知：记录看板操作日志，并展示与当前用户相关的通知。
- 看板统计视图：各列任务数量、成员任务分布饼图、本周完成任务趋势折线图、逾期任务列表。

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 认证 | JWT（github.com/golang-jwt/jwt/v5）+ 管理员/编辑者/观察者 RBAC |
| 前端 | React 18 + TypeScript + Ant Design |
| 构建 | Vite |
| 拖拽 | @hello-pangea/dnd |
| 状态管理 | Zustand |
| HTTP 客户端 | Axios |
| 其他前端依赖 | dayjs、lodash |

## 快速启动（Docker Compose 一键部署）

```bash
cp .env.example .env
docker compose --env-file .env up -d --build
```

启动完成后：

- 前端访问地址：http://127.0.0.1:18503
- 后端健康检查：http://127.0.0.1:19503/healthz
- 数据库地址：127.0.0.1:57503

首次启动会由后端 GORM 自动迁移数据库表结构。

## 本地开发

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

前端：

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器会把 `/api` 与 `/uploads` 代理到 `http://127.0.0.1:19503`。

## 访问地址

| 服务 | 地址 |
| --- | --- |
| 前端 | http://127.0.0.1:18503 |
| 后端 API | http://127.0.0.1:19503/api/v1 |
| 健康检查 | http://127.0.0.1:19503/healthz |
| MySQL | 127.0.0.1:57503 |

## 项目目录结构

```text
.
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── handler/
│   │   ├── router/
│   │   ├── middleware/
│   │   ├── dto/
│   │   ├── constants/
│   │   └── response/
│   ├── migrations/
│   ├── api/
│   ├── deploy/
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── store/
│   │   └── types/
│   ├── nginx.conf
│   └── Dockerfile
├── database/init.sql
├── docker-compose.yml
├── .env.example
└── README.md
```

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `gbkanban` | Docker Compose 项目名 |
| `DB_NAME` | `gbkanban` | 数据库名 |
| `DB_USER` | `gbkanban` | 数据库用户 |
| `DB_PASSWORD` | `gbkanban` | 数据库密码 |
| `DB_ROOT_PASSWORD` | `gbkanban_root` | MySQL root 密码 |
| `JWT_SECRET` | 开发默认值 | JWT 签名密钥，生产环境务必修改 |
| `FRONTEND_PORT` | `18503` | 前端宿主机端口 |
| `BACKEND_PORT` | `19503` | 后端宿主机端口 |
| `DB_PORT` | `57503` | 数据库宿主机端口 |

## Docker 部署说明

- `docker-compose.yml` 不使用 `version` 字段，顶层声明 `name: gbkanban`。
- 数据库数据通过命名卷 `gbkanban_mysql_data` 持久化，附件通过 `gbkanban_uploads_data` 持久化。
- 后端镜像采用 Go 多阶段构建，前端镜像采用 Node + Nginx 多阶段构建。
- 前端 Nginx 将 `/api/` 反向代理到 `http://backend:8080/`，并支持 SPA 路由回退到 `index.html`。

## API 清单

所有接口统一返回 `{ "code": 0, "message": "ok", "data": ... }`。

### 认证

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/me`

### 工作区与成员

- `GET /api/v1/workspaces`
- `POST /api/v1/workspaces`
- `GET /api/v1/workspaces/:workspace_id`
- `PUT /api/v1/workspaces/:workspace_id`
- `DELETE /api/v1/workspaces/:workspace_id`
- `GET /api/v1/workspaces/:workspace_id/members`
- `POST /api/v1/workspaces/:workspace_id/members`
- `PUT /api/v1/workspaces/:workspace_id/members/:user_id`
- `DELETE /api/v1/workspaces/:workspace_id/members/:user_id`

### 看板与列

- `GET /api/v1/workspaces/:workspace_id/boards`
- `POST /api/v1/workspaces/:workspace_id/boards`
- `GET /api/v1/workspaces/:workspace_id/boards/:board_id`
- `PUT /api/v1/workspaces/:workspace_id/boards/:board_id`
- `DELETE /api/v1/workspaces/:workspace_id/boards/:board_id`
- `GET /api/v1/workspaces/:workspace_id/boards/:board_id/columns`
- `POST /api/v1/workspaces/:workspace_id/boards/:board_id/columns`
- `PUT /api/v1/workspaces/:workspace_id/columns/:column_id`
- `DELETE /api/v1/workspaces/:workspace_id/columns/:column_id`

### 任务与详情

- `GET /api/v1/workspaces/:workspace_id/boards/:board_id/tasks`
- `POST /api/v1/workspaces/:workspace_id/boards/:board_id/tasks`
- `GET /api/v1/workspaces/:workspace_id/tasks/:task_id`
- `PUT /api/v1/workspaces/:workspace_id/tasks/:task_id`
- `DELETE /api/v1/workspaces/:workspace_id/tasks/:task_id`
- `PATCH /api/v1/workspaces/:workspace_id/tasks/:task_id/move`
- `GET /api/v1/workspaces/:workspace_id/tags`
- `POST /api/v1/workspaces/:workspace_id/tags`
- `DELETE /api/v1/workspaces/:workspace_id/tags/:tag_id`
- `POST /api/v1/workspaces/:workspace_id/tasks/:task_id/subtasks`
- `PUT /api/v1/workspaces/:workspace_id/subtasks/:subtask_id`
- `DELETE /api/v1/workspaces/:workspace_id/subtasks/:subtask_id`
- `POST /api/v1/workspaces/:workspace_id/tasks/:task_id/comments`
- `POST /api/v1/workspaces/:workspace_id/tasks/:task_id/attachments`
- `DELETE /api/v1/workspaces/:workspace_id/attachments/:attachment_id`

### 活动、统计与通知

- `GET /api/v1/workspaces/:workspace_id/activities`
- `GET /api/v1/workspaces/:workspace_id/boards/:board_id/activities`
- `GET /api/v1/workspaces/:workspace_id/stats`
- `GET /api/v1/notifications`
- `GET /api/v1/notifications/unread-count`
- `PUT /api/v1/notifications/:notification_id/read`
- `PUT /api/v1/notifications/read-all`

## License

MIT

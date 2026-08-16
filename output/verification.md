# 团队任务看板 验证报告

## 项目信息

- 项目目录：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/综合实用主题项目提示词/cy-303`
- 提示词文件：`cy-303.md`
- 提示词 SHA-256：
  - 开始：`4e1932730e80753e49266342fe6daffa00d478949ba9577aa4a06fccf1d0167e`
  - 结束：`4e1932730e80753e49266342fe6daffa00d478949ba9577aa4a06fccf1d0167e`
  - 前后一致，提示词文件未被修改。

## 生成的主要目录与文件

```text
.
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── model/*.go
│   │   ├── repository/*.go
│   │   ├── service/*.go
│   │   ├── handler/*.go
│   │   ├── router/router.go
│   │   ├── middleware/*.go
│   │   ├── dto/*.go
│   │   ├── constants/*.go
│   │   └── response/response.go
│   ├── migrations/000001_init.up.sql
│   ├── api/openapi.yaml
│   ├── deploy/docker-compose.backend.yml
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── store/
│   │   └── types/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── database/init.sql
├── docker-compose.yml
├── .env.example
├── .gitignore
├── README.md
└── output/
    ├── verification.md
    ├── api-test.log
    └── screenshots/*.png
```

## 构建验证

### 后端

在 `backend/` 下执行：

- `go mod tidy`：通过
- `go vet ./...`：通过
- `go build ./...`：通过
- `go test ./...`：通过（service 与 repository 均包含表驱动测试）

依赖版本已固定在 Go 1.22 可构建的范围，`go.mod` 声明 `go 1.22`，Docker 使用 `golang:1.22-alpine` 构建成功。

### 前端

在 `frontend/` 下执行：

- `npm install`：通过
- `npm run build`（含 `tsc -b` 类型检查 + Vite 构建）：通过

### Docker Compose

- `docker compose --env-file .env config --quiet`：通过
- `docker compose --env-file .env up -d --build --wait`：db、backend、frontend 三个服务全部 `healthy`。
- 端口映射：
  - 前端 `18503 -> 80`
  - 后端 `19503 -> 8080`
  - 数据库 `57503 -> 3306`

## 接口实测

主要流程已通过 `curl` 全量验证，完整日志见 `output/api-test.log`：

- `/healthz` 与 `/health`：返回正常。
- 注册/登录/`/me`：`alice`、`bob`、`carol` 注册与登录成功。
- 工作区/看板/列：创建工作区、创建看板（自动生成 To Do / In Progress / Done）、新增自定义列「Review」成功。
- 标签与任务 CRUD：创建标签、创建任务、任务详情字段（优先级/负责人/截止日期/标签）返回正确。
- 子任务/评论/附件：新增子任务、勾选完成、添加评论、上传附件均成功。
- 拖拽/状态流转：`PATCH /tasks/:id/move` 将任务从 To Do 移到 Done 成功，排序更新正确。
- 筛选与搜索：按优先级、关键词、标签筛选均返回正确结果。
- 成员与角色：邀请 bob（editor）、carol（viewer），更新角色成功。
- RBAC：viewer 创建任务返回 `403`；editor 创建任务成功。
- 活动日志：记录 task_created、task_moved、task_updated 等操作。
- 通知：将任务分配给 bob 后，bob 收到 `assigned` 通知，标记已读后未读数归零。
- 统计：列任务数、成员分布、本周趋势、逾期任务列表返回正常。

## 前端浏览器验证（应用内浏览器）

使用内置 Browser 插件（应用内浏览器）访问 `http://127.0.0.1:18503`，页面均真实渲染：

- 登录页：输入框与登录按钮渲染正常，`alice` 登录成功。
- 工作区页：显示「产品研发」工作区与所有者信息。
- 看板页：To Do / In Progress / Done / Review 四列渲染，任务卡片显示标题、优先级、标签、负责人头像与截止日期。
- 任务详情：点击卡片弹出侧边栏，显示详情、子任务、评论、附件四个页签，Markdown 描述预览正常。
- 拖拽：通过浏览器模拟拖拽，将「联调拖拽」任务从 To Do 移动到 In Progress，前端实时更新，后端 `column_id` 同步变更。
- 统计页：显示各列任务数、成员任务分布、本周完成任务趋势、逾期任务列表。
- `/api` 代理：前端 `/` 与 SPA 路由均返回 200，`/api/v1/...` 经 Nginx 正确转发到后端并返回业务响应。

截图已保存到 `output/screenshots/`：

- `gbkanban-01-login.png`
- `gbkanban-02-workspaces.png`
- `gbkanban-03-board.png`
- `gbkanban-04-task-detail.png`
- `gbkanban-05-stats.png`
- `gbkanban-06-dragged.png`

## 清理验证

执行 `docker compose --env-file .env down -v --remove-orphans` 后：

- 容器、网络、命名卷均已删除。
- `lsof` 检查宿主机端口 `18503`、`19503`、`57503` 均无监听。
- `docker compose --env-file .env ps -a` 无残留。

## Git 初始化

- `git init` 已完成。
- 本地用户：`blueship581 <brysj.hhrhl.g@gmail.com>`。
- 提交信息：`feat: 团队任务看板`。
- 工作区 `git status --short` 干净。

## 已知问题

1. 生产构建产物为单个 JS chunk，体积约 1.5 MB（gzip 约 498 KB），Vite 提示可进一步做代码分割；不影响功能。
2. 统计页的饼图与折线图使用原生 SVG 实现，未引入第三方图表库，视觉上能满足需求但交互能力有限。
3. 附件删除仅删除数据库记录，未同步清理卷中的物理文件（已通过命名卷持久化，不影响演示流程）。
4. `database/init.sql` 为数据库/用户 bootstrap，实际表结构由后端 GORM AutoMigrate 在启动时创建。

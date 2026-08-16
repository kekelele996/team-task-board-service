# 团队任务看板（cy-303）审计报告

## 项目信息
- 项目目录：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/综合实用主题项目提示词/cy-303`
- 提示词：`cy-303.md`
- 技术栈：React 18 + TypeScript + Ant Design + Vite + @hello-pangea/dnd + Zustand；Go 1.22 + Gin + GORM + MySQL 8
- 端口：前端 18503，后端 19503，数据库 57503，服务名 gbkanban

## 1. 静态验证
- `go vet ./...`：通过
- `go build ./...`：通过
- `go test ./...`：通过
- `npm run build`：通过（仅提示 chunk 体积警告）
- `docker compose --env-file .env config --quiet`：通过

## 2. 运行验证
- `docker compose --env-file .env up -d --build --wait`：db/backend/frontend 均 healthy
- `GET /health`、`GET /healthz`：200
- 前端 `/` 及 Nginx `/api/` 代理：正常

### API 流程（后端直连 19503 与 Nginx 代理 18503）
- 注册 / 登录 / me
- 工作区、看板创建与默认三列，自定义列
- 标签创建与列表
- 成员邀请、成员列表、角色更新（admin/editor/viewer RBAC）
- 任务创建、列表、详情、更新、删除
- 拖拽移动与同列排序（`PATCH /tasks/:id/move`）
- 子任务、评论、附件上传
- 按关键词、负责人、优先级、标签、截止日期筛选
- 活动日志（看板 / 工作区）
- 通知（列表、未读数、全部已读）
- 统计（列数量、成员分布、周趋势、逾期列表）

证据：
- `output/api-test.log`（主流程）
- `output/api-test-fixed.log`（修复后的删除与 RBAC 回归）
- 浏览器截图：`output/screenshots/gbkanban-*-new.png`

## 3. 前端浏览器验证（In-app Browser，非外部 Chrome）
- 登录：输入邮箱/密码成功进入工作区列表
- 看板：正确展示三列及任务卡片
- 任务详情：标题、描述、负责人、优先级、标签、子任务/评论/附件 Tab 正常
- 拖拽：将「Write tests」从 To Do 拖到 Done，UI 与后端状态均更新，提示「任务已移动」
- 统计：列数量、成员任务分布、本周完成任务趋势正常显示

## 4. 发现的问题（已修复）
1. 删除带有标签/子任务/评论/附件的任务返回 500。根因：`task_repo.Delete` 直接删除 `tasks`，触发 `task_tags` 等外键约束。
2. 删除看板返回 500。根因：`board_repo.Delete` 未先删除列与任务，触发 `board_columns` 外键约束。
3. 删除列（列下存在任务）返回 500。根因：`column_repo.Delete` 未处理任务及任务子表。
4. 删除工作区返回 500。根因：`workspace_repo.Delete` 未先清理成员、看板、列、任务等依赖数据。
5. 任务创建 / 移动未校验 `column_id` 是否属于目标看板，可跨看板传列 ID 导致数据不一致。
6. 看板列列表、看板活动日志接口只校验成员身份，未校验 `board_id` 是否属于路径中的 `workspace_id`，存在跨工作区越权读取风险。

## 5. 修复内容
- `backend/internal/repository/task_repo.go`：事务内先删除 `task_tags`、`subtasks`、`comments`、`attachments`，再删除任务。
- `backend/internal/repository/board_repo.go`：事务内按 `board_id` 清理任务子表、任务、列，再删除看板。
- `backend/internal/repository/column_repo.go`：事务内按 `column_id` 清理任务子表、任务，再删除列。
- `backend/internal/repository/workspace_repo.go`：事务内按 `workspace_id` 清理任务子表、任务、列、看板、成员、活动日志，再删除工作区。
- `backend/internal/service/task_service.go`：创建任务与移动任务时校验列属于目标看板；`ListBoardActivities` 增加 workspace/board 归属校验。
- `backend/internal/handler/board_handler.go`、`task_handler.go`：列列表和看板活动日志增加 `workspace_id` 归属解析与校验。

## 6. 修复后回归
- 删除任务：`code:0`
- 删除有任务的列：`code:0`
- 删除有任务的看板：`code:0`
- 删除工作区：`code:0`
- 跨看板列 ID 创建/移动任务：返回 `40400 resource not found`
- 全部 `go vet` / `go build` / `go test` / `npm run build` / `docker compose config --quiet` 仍通过

## 7. 最终状态
- 项目可正常构建、测试、部署，主业务流程与 RBAC 均可用。
- 已发现并修复的均为真实功能/权限缺陷，未改动 `cy-303.md` 提示词。

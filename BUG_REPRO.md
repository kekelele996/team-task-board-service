# BUG_REPRO

## Bug 是什么
任务初始位置计算错位，首条任务落不到 0 位；任务列表排序键从位置+id 改成只按 id；单任务移动时没更新所属列；移动请求丢了位置非负校验。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestTaskMoveAndReorder 失败：新建任务位置不为 0、重排后顺序不变、跨列移动后仍留在原列。

# BUG_REPRO

## Bug 是什么
任务和看板查询在未命中时返回 nil,nil，上层拿到 nil 后直接解引用，导致空指针 panic。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestMissingBoardAndTaskReturnNotFound 失败：查不存在的看板时发生 invalid memory address or nil pointer dereference。
- repository.TestTaskRepo_ListFilterAndMove 失败：删除后 FindByID 返回 nil 错误。

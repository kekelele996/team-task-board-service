# BUG_REPRO

## Bug 是什么
用户、任务、成员查询在未命中时把哨兵错误用 %v 包了一层，errors.Is 全部失效；登录时又丢掉了未命中转 401 的分支。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestAuthErrorChain 失败：新邮箱注册报错、未知邮箱登录返回服务器错误。
- repository.TestTaskRepo_ListFilterAndMove 失败：删除后 FindByID 不再命中 ErrNotFound。
- service.TestAuthService_Login/unknown_user 失败：未返回 invalid credentials。

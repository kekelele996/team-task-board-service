# BUG_REPRO

## Bug 是什么
默认列名顺序被写反，创建看板时三列位置全部写死成 0，列表查询又按位置倒序排，导致新建看板默认列顺序错误。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestBoardDefaultColumnsOrder 失败：实际得到 [Done In Progress To Do]，期望 [To Do In Progress Done]。

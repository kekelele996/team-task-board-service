# BUG_REPRO

## Bug 是什么
评论任务时把通知重复写了两条；通知里操作人写成了接收人自己；未读统计去掉了用户过滤；通知默认已读状态写反。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestCommentNotificationChain 失败：接收人拿到 2 条通知、未读数为 0、通知操作人不是评论人。

# FAQ - 常见问题

## 安装和配置

### Q: 如何获取 Apple App Store Connect API 凭证？

A: 
1. 登录 [App Store Connect](https://appstoreconnect.apple.com/)
2. 进入 "用户和访问" → "密钥"
3. 点击 "+" 创建新密钥
4. 下载 `.p8` 文件（只能下载一次，请妥善保管）
5. 记录 Issuer ID 和 Key ID

### Q: 如何获取 Google Play Developer API 凭证？

A:
1. 登录 [Google Cloud Console](https://console.cloud.google.com/)
2. 创建或选择项目
3. 启用 "Google Play Android Developer API"
4. 创建服务账号并下载 JSON 密钥文件
5. 在 [Google Play Console](https://play.google.com/console/) 中授予服务账号权限

### Q: 支持哪些 Go 版本？

A: 本 SDK 支持 Go 1.21 及以上版本。

## 使用问题

### Q: 为什么获取不到评论？

A: 可能的原因：
1. 凭证配置错误
2. 时间窗口内没有新评论
3. API 权限不足
4. 网络连接问题

建议：
- 检查错误信息
- 尝试扩大时间范围（增加 `Since` 参数）
- 验证凭证是否正确

### Q: 如何只获取低分评论？

A: 目前 SDK 不直接支持按评分过滤，但您可以：

```go
reviews, _ := client.FetchAppleReviews(ctx, nil)
for _, r := range reviews {
    rating, _ := strconv.Atoi(r.ReviewRating)
    if rating <= 2 {
        // 处理低分评论
    }
}
```

### Q: 可以同时获取多个应用的评论吗？

A: 需要为每个应用创建单独的 Client：

```go
client1, _ := appstore.NewAppleOnly(appstore.AppleConfig{
    AppID: "app-id-1",
    // ...
})

client2, _ := appstore.NewAppleOnly(appstore.AppleConfig{
    AppID: "app-id-2",
    // ...
})
```

### Q: 如何处理 API 限流？

A: 
1. 使用合理的时间间隔
2. 实现指数退避重试
3. 使用 `OnPage` 回调控制拉取速度

```go
opts := &appstore.FetchAppleOptions{
    OnPage: func(items []appstore.ReviewItem) bool {
        time.Sleep(1 * time.Second) // 延迟
        return true
    },
}
```

## 错误处理

### Q: 遇到 "apple configuration not set" 错误怎么办？

A: 这表示 Apple 配置未正确设置。检查：
1. 是否调用了 `New()` 或 `NewAppleOnly()`
2. 配置参数是否完整
3. 是否通过了 `Validate()` 验证

### Q: 遇到 JWT 签名错误怎么办？

A: 可能的原因：
1. `.p8` 文件内容不完整或损坏
2. Key ID 或 Issuer ID 错误
3. 私钥格式不正确

建议重新下载 `.p8` 文件并检查配置。

### Q: Google API 返回 403 错误？

A: 检查：
1. 服务账号是否在 Google Play Console 中授权
2. API 是否已启用
3. 包名是否正确

## 性能优化

### Q: 如何提高获取评论的速度？

A:
1. 使用合理的 `PerPage` 参数（Apple 最大 200）
2. 缩小时间范围
3. 使用 `Limit` 参数限制数量
4. 考虑并发获取多个应用的评论

### Q: 可以缓存评论数据吗？

A: 可以，但需要注意：
1. 评论可能被用户修改或删除
2. 建议定期刷新缓存
3. 注意数据隐私合规

## 开发和测试

### Q: 如何在不调用真实 API 的情况下测试？

A: 可以：
1. 使用 mock 数据
2. 创建测试用的 interface
3. 使用 httptest 模拟 API 响应

### Q: 有示例项目吗？

A: 查看 [examples/](../examples/) 目录，包含多个使用示例。

## 其他

### Q: 支持其他应用商店吗（如华为、小米）？

A: 目前仅支持 Apple App Store 和 Google Play。如需支持其他平台，欢迎提交 PR。

### Q: 可以商业使用吗？

A: 可以，本项目使用 MIT License，允许商业使用。

### Q: 如何贡献代码？

A: 请查看 [CONTRIBUTING.md](../CONTRIBUTING.md)。

---

如果您的问题未在此列出，欢迎：
- 创建 [GitHub Issue](https://github.com/Leslack-hub/appstore-reviews-go/issues)
- 参与 [Discussions](https://github.com/Leslack-hub/appstore-reviews-go/discussions)

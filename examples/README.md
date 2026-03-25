# Examples

本目录包含 appstore-reviews-go SDK 的使用示例。

## 示例列表

### [basic](basic/main.go)
基础用法示例，展示如何：
- 创建双平台客户端
- 获取 Apple 和 Google 评论
- 提交回复

### [advanced](advanced/main.go)
高级用法示例，展示如何：
- 创建单平台客户端
- 使用自定义查询选项
- 设置时间范围和数量限制
- 使用分页回调
- 过滤特定评论

## 运行示例

在运行示例之前，您需要：

1. 准备 Apple 凭证（如果使用 Apple 功能）：
   - 从 App Store Connect 下载 `.p8` 私钥文件
   - 获取 Issuer ID 和 Key ID
   - 获取您的 App ID

2. 准备 Google 凭证（如果使用 Google 功能）：
   - 从 Google Cloud Console 创建服务账号
   - 下载服务账号 JSON 文件
   - 在 Google Play Console 授予权限

3. 修改示例代码中的配置信息

4. 运行示例：

```bash
# 运行基础示例
go run examples/basic/main.go

# 运行高级示例
go run examples/advanced/main.go
```

## 注意事项

- 请勿将凭证文件提交到版本控制系统
- 建议使用环境变量或配置文件管理敏感信息
- 测试时注意 API 调用频率限制

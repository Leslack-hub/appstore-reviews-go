# Architecture

本文档描述 appstore-reviews-go SDK 的架构设计。

## 设计原则

1. **简单易用**：提供清晰的 API，降低使用门槛
2. **灵活可扩展**：支持多种配置方式和查询选项
3. **错误明确**：提供清晰的错误类型和错误信息
4. **平台独立**：支持单平台或多平台使用
5. **上下文支持**：所有操作支持 context 控制

## 核心组件

### Client

`Client` 是 SDK 的核心，负责协调 Apple 和 Google 两个平台的操作。

```
┌─────────────────────────────────────┐
│           Client                    │
├─────────────────────────────────────┤
│ - apple: AppleConfig                │
│ - google: *googleClient             │
│ - appleEnabled: bool                │
├─────────────────────────────────────┤
│ + FetchAppleReviews()               │
│ + FetchGoogleReviews()              │
│ + ReplyAppleReview()                │
│ + ReplyGoogleReview()               │
│ + IsAppleEnabled()                  │
│ + IsGoogleEnabled()                 │
└─────────────────────────────────────┘
```

### Apple 实现 (apple.go)

Apple 平台的实现直接在 `apple.go` 中，主要功能：

- JWT Token 生成（使用 ECDSA 签名）
- App Store Connect API 调用
- 评论数据转换为统一格式

```
generateAppleToken() → JWT Token
         ↓
fetchAppleReviews() → HTTP Request → Apple API
         ↓
appleReviewResponse → []ReviewItem
```

### Google 实现 (google.go)

Google 平台使用 `googleClient` 封装，主要功能：

- OAuth2 认证
- Google Play Developer API 调用
- 评论数据转换为统一格式

```
newGoogleClient() → googleClient
         ↓
fetchReviews() → Google API Client → Google API
         ↓
androidpublisher.Review → []ReviewItem
```

## 数据流

### 获取评论流程

```
User Code
    ↓
Client.FetchAppleReviews(ctx, opts)
    ↓
fetchAppleReviews(ctx, config, opts)
    ↓
generateAppleToken() → JWT
    ↓
HTTP GET → Apple API
    ↓
Parse JSON → appleReviewResponse
    ↓
Transform → []ReviewItem
    ↓
Return to User
```

### 提交回复流程

```
User Code
    ↓
Client.ReplyAppleReview(ctx, reviewID, content)
    ↓
replyAppleReview(ctx, config, reviewID, content)
    ↓
generateAppleToken() → JWT
    ↓
HTTP POST → Apple API
    ↓
Check Status Code
    ↓
Return error or nil
```

## 错误处理

SDK 使用分层错误处理策略：

1. **配置错误**：在创建 Client 时验证
   - `ErrInvalidConfig`
   
2. **运行时错误**：在调用方法时检查
   - `ErrAppleNotConfigured`
   - `ErrGoogleNotConfigured`
   
3. **API 错误**：包装底层 API 错误
   - HTTP 错误
   - JSON 解析错误
   - 认证错误

## 扩展性

### 添加新平台

要添加新平台支持（如华为应用市场），需要：

1. 创建新文件 `huawei.go`
2. 实现 `huaweiClient` 结构
3. 在 `Client` 中添加相应字段和方法
4. 在 `model.go` 中添加平台常量

### 添加新功能

要添加新功能（如获取评论统计），需要：

1. 在相应平台文件中实现底层逻辑
2. 在 `Client` 中添加公开方法
3. 更新文档和示例

## 性能考虑

### 分页处理

- Apple API 使用 cursor-based 分页
- Google API 使用 token-based 分页
- 支持 `OnPage` 回调实现流式处理

### 并发安全

- `Client` 是并发安全的
- 可以在多个 goroutine 中共享同一个 `Client` 实例

### 资源管理

- 使用 `context.Context` 控制请求生命周期
- HTTP Client 设置合理的超时时间
- 支持取消操作

## 测试策略

1. **单元测试**：测试配置验证、错误处理等
2. **集成测试**：测试与真实 API 的交互（需要凭证）
3. **Mock 测试**：使用 mock server 测试 HTTP 交互

## 依赖管理

核心依赖：
- `github.com/golang-jwt/jwt/v4`：JWT 生成
- `golang.org/x/oauth2`：OAuth2 认证
- `google.golang.org/api`：Google API 客户端

最小化依赖原则：
- 不引入不必要的第三方库
- 优先使用标准库
- 定期更新依赖版本

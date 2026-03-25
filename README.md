# appstore-reviews-go

[![Go Reference](https://pkg.go.dev/badge/github.com/Leslack-hub/appstore-reviews-go.svg)](https://pkg.go.dev/github.com/Leslack-hub/appstore-reviews-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Leslack-hub/appstore-reviews-go)](https://goreportcard.com/report/github.com/Leslack-hub/appstore-reviews-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CI](https://github.com/Leslack-hub/appstore-reviews-go/workflows/CI/badge.svg)](https://github.com/Leslack-hub/appstore-reviews-go/actions)

Go SDK，用于从 Apple App Store 和 Google Play 获取用户评论，以及提交开发者回复。

## 特性

- ✅ 支持 Apple App Store Connect API
- ✅ 支持 Google Play Developer API
- ✅ 统一的数据结构和接口
- ✅ 灵活的查询选项（时间范围、数量限制、分页等）
- ✅ 完善的错误处理
- ✅ 支持单平台或双平台配置
- ✅ 上下文支持，可取消操作
- ✅ 并发安全
- ✅ 完整的测试覆盖

## 目录

- [安装](#安装)
- [快速上手](#快速上手)
- [API 参考](#api-参考)
- [示例代码](#示例代码)
- [文档](#文档)
- [开发](#开发)
- [贡献](#贡献)
- [License](#license)

## 安装

```bash
go get github.com/Leslack-hub/appstore-reviews-go
```

要求：Go 1.21+

## 快速上手

### 基础用法

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    appstore "github.com/Leslack-hub/appstore-reviews-go"
)

func main() {
    certPEM, _ := os.ReadFile("AuthKey_XXXXXXXXXX.p8")
    credJSON, _ := os.ReadFile("google-service-account.json")

    // 创建双平台客户端
    client, err := appstore.New(context.Background(),
        appstore.AppleConfig{
            IssuerID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
            KeyID:    "XXXXXXXXXX",
            Cert:     string(certPEM),
            AppID:    "6737681541",
        },
        appstore.GoogleConfig{
            CredentialsJSON: credJSON,
            PackageName:     "com.example.myapp",
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // 获取苹果评论（默认最近 48 小时）
    appleReviews, err := client.FetchAppleReviews(ctx, nil)
    if err != nil {
        log.Printf("获取苹果评论失败: %v", err)
    }
    fmt.Printf("苹果评论数: %d\n", len(appleReviews))

    // 获取谷歌评论（默认最近 2 小时）
    googleReviews, err := client.FetchGoogleReviews(ctx, nil)
    if err != nil {
        log.Printf("获取谷歌评论失败: %v", err)
    }
    fmt.Printf("谷歌评论数: %d\n", len(googleReviews))

    // 回复评论
    if len(appleReviews) > 0 {
        err := client.ReplyAppleReview(ctx, appleReviews[0].ReviewId, "感谢您的反馈！")
        if err != nil {
            log.Printf("回复失败: %v", err)
        }
    }
}
```

### 单平台客户端

```go
// 仅使用 Apple
client, err := appstore.NewAppleOnly(appstore.AppleConfig{
    IssuerID: "your-issuer-id",
    KeyID:    "your-key-id",
    Cert:     certPEM,
    AppID:    "your-app-id",
})

// 仅使用 Google
client, err := appstore.NewGoogleOnly(ctx, appstore.GoogleConfig{
    CredentialsJSON: credJSON,
    PackageName:     "com.example.app",
})
```

### 高级用法

```go
import "time"

// 自定义查询选项
reviews, err := client.FetchAppleReviews(ctx, &appstore.FetchAppleOptions{
    Limit:   10,                 // 最多获取 10 条
    Since:   7 * 24 * time.Hour, // 最近 7 天
    PerPage: 50,                 // 每页 50 条
    Sort:    "-createdDate",     // 按创建时间倒序
    OnPage: func(items []appstore.ReviewItem) bool {
        fmt.Printf("收到 %d 条评论\n", len(items))
        return true // 返回 false 可中止拉取
    },
})
```

## API 参考

### 客户端创建

```go
// 双平台客户端
func New(ctx context.Context, apple AppleConfig, google GoogleConfig) (*Client, error)

// 单平台客户端
func NewAppleOnly(apple AppleConfig) (*Client, error)
func NewGoogleOnly(ctx context.Context, google GoogleConfig) (*Client, error)
```

### 配置结构

#### `AppleConfig`

| 字段 | 类型 | 说明 |
|------|------|------|
| `IssuerID` | `string` | App Store Connect Issuer ID |
| `KeyID` | `string` | API Key ID |
| `Cert` | `string` | ECDSA 私钥 PEM 内容（.p8 文件） |
| `AppID` | `string` | 苹果应用 ID（数字字符串） |

#### `GoogleConfig`

| 字段 | 类型 | 说明 |
|------|------|------|
| `CredentialsJSON` | `[]byte` | 服务账号 JSON 文件内容 |
| `PackageName` | `string` | 安卓应用包名 |

### 查询选项

#### `FetchAppleOptions`

| 字段 | 类型 | 说明 |
|------|------|------|
| `Limit` | `int` | 最大评论数，0 表示不限制 |
| `Since` | `time.Duration` | 时间窗口，0 表示不限制 |
| `PerPage` | `int` | 每页数量（1-200），默认 200 |
| `Sort` | `string` | 排序方式，默认 "-createdDate" |
| `QueryParams` | `url.Values` | 额外查询参数 |
| `OnPage` | `func([]ReviewItem) bool` | 分页回调 |

#### `FetchGoogleOptions`

| 字段 | 类型 | 说明 |
|------|------|------|
| `Limit` | `int` | 最大评论数，0 表示不限制 |
| `Since` | `time.Duration` | 时间窗口，0 表示不限制 |
| `TranslationLanguage` | `string` | 翻译语言 |
| `OnPage` | `func([]ReviewItem) bool` | 分页回调 |

### 响应结构

#### `ReviewItem`

| 字段 | 类型 | 说明 |
|------|------|------|
| `Platform` | `uint8` | `1`=Apple, `2`=Google |
| `ReviewId` | `string` | 评论唯一 ID |
| `ReviewTitle` | `string` | 评论标题 |
| `OriginalContent` | `string` | 原始评论内容 |
| `TranslatedContent` | `string` | 翻译内容（仅 Google Play） |
| `ReviewNickname` | `string` | 评论者昵称 |
| `ReviewRating` | `string` | 星级（"1"~"5"） |
| `ReviewLanguage` | `string` | 语言/地区代码 |
| `ReviewExtra` | `map[string]any` | 附加信息 |
| `CreatedAt` | `string` | 评论时间（RFC3339） |

### 客户端方法

```go
// 获取评论
func (c *Client) FetchAppleReviews(ctx context.Context, opts *FetchAppleOptions) ([]ReviewItem, error)
func (c *Client) FetchGoogleReviews(ctx context.Context, opts *FetchGoogleOptions) ([]ReviewItem, error)

// 回复评论
func (c *Client) ReplyAppleReview(ctx context.Context, reviewID, content string) error
func (c *Client) ReplyGoogleReview(ctx context.Context, reviewID, content string) error

// 状态检查
func (c *Client) IsAppleEnabled() bool
func (c *Client) IsGoogleEnabled() bool
```

### 错误类型

```go
var (
    ErrAppleNotConfigured  = errors.New("apple configuration not set")
    ErrGoogleNotConfigured = errors.New("google configuration not set")
    ErrInvalidConfig       = errors.New("invalid configuration")
)
```

### 平台常量

```go
const (
    PlatformApple  uint8 = 1  // Apple App Store
    PlatformGoogle uint8 = 2  // Google Play
)
```

## 示例代码

查看 [examples/](examples/) 目录获取更多示例：

- [basic](examples/basic/main.go) - 基础用法
- [advanced](examples/advanced/main.go) - 高级选项和自定义查询
- [single-platform](examples/single-platform/main.go) - 单平台客户端
- [error-handling](examples/error-handling/main.go) - 错误处理最佳实践

## 文档

- [API 文档](https://pkg.go.dev/github.com/Leslack-hub/appstore-reviews-go) - 完整的 API 参考
- [架构设计](doc/ARCHITECTURE.md) - 了解 SDK 内部设计
- [常见问题](doc/FAQ.md) - 常见问题解答
- [安全政策](SECURITY.md) - 安全最佳实践

## 开发

```bash
# 克隆仓库
git clone https://github.com/Leslack-hub/appstore-reviews-go.git
cd appstore-reviews-go

# 安装依赖
go mod download

# 运行测试
make test

# 代码检查
make lint

# 格式化代码
make fmt

# 查看所有命令
make help
```

## 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

在提交 PR 之前，请确保：
- 代码通过所有测试
- 添加了必要的测试
- 更新了相关文档
- 遵循 Go 代码规范

## 变更日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本历史。

## 相关项目

- [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi)
- [Google Play Developer API](https://developers.google.com/android-publisher)

## License

MIT License - see [LICENSE](LICENSE) for details

---

如有问题或建议，欢迎：
- 创建 [Issue](https://github.com/Leslack-hub/appstore-reviews-go/issues)
- 提交 [Pull Request](https://github.com/Leslack-hub/appstore-reviews-go/pulls)
- 参与 [Discussions](https://github.com/Leslack-hub/appstore-reviews-go/discussions)

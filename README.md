# appstore-reviews-go

Go SDK，用于从 Apple App Store 和 Google Play 获取用户评论，以及提交开发者回复。

## 安装

```bash
go get github.com/Leslack-hub/appstore-reviews-go
```

## 快速上手

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Leslack-hub/appstore-reviews-go"
)

func main() {
    // 读取苹果 .p8 私钥文件
    certPEM, err := os.ReadFile("AuthKey_XXXXXXXXXX.p8")
    if err != nil {
        log.Fatal(err)
    }

    // 读取谷歌服务账号 JSON 凭证文件
    credJSON, err := os.ReadFile("google-service-account.json")
    if err != nil {
        log.Fatal(err)
    }

    // 创建 SDK Client
    client, err := appstore.New(context.Background(),
        appstore.AppleConfig{
            IssuerID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", // App Store Connect Issuer ID
            KeyID:    "XXXXXXXXXX",                           // API Key ID
            Cert:     string(certPEM),                        // .p8 私钥内容
            AppID:    "6737681541",                           // 苹果 App ID
        },
        appstore.GoogleConfig{
            CredentialsJSON: credJSON,              // 服务账号 JSON 内容
            PackageName:     "com.example.myapp",  // 安卓包名
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // ── 获取苹果评论 ──────────────────────────────────────────
    // limit=0: 拉取 48h 内全部新评论
    // limit=10: 最多拉取 10 条最新评论
    appleReviews, err := client.FetchAppleReviews(ctx, 0, 0)
    if err != nil {
        log.Printf("获取苹果评论失败: %v", err)
    }
    fmt.Printf("苹果评论数: %d\n", len(appleReviews))
    for _, r := range appleReviews {
        fmt.Printf("  [%s] %s ★%s: %s\n", r.ReviewLanguage, r.ReviewNickname, r.ReviewRating, r.OriginalContent)
    }

    // ── 获取谷歌评论 ──────────────────────────────────────────
    // limit=0: 拉取 2h 内全部新评论
    googleReviews, err := client.FetchGoogleReviews(ctx, 0, 0)
    if err != nil {
        log.Printf("获取谷歌评论失败: %v", err)
    }
    fmt.Printf("谷歌评论数: %d\n", len(googleReviews))

    // ── 提交苹果回复 ──────────────────────────────────────────
    if len(appleReviews) > 0 {
        r := appleReviews[0]
        if err := client.ReplyAppleReview(ctx, r.ReviewId, "感谢您的反馈！"); err != nil {
            log.Printf("回复苹果评论失败: %v", err)
        } else {
            fmt.Println("苹果回复提交成功")
        }
    }

    // ── 提交谷歌回复 ──────────────────────────────────────────
    if len(googleReviews) > 0 {
        r := googleReviews[0]
        if err := client.ReplyGoogleReview(ctx, r.ReviewId, "感谢您的反馈！"); err != nil {
            log.Printf("回复谷歌评论失败: %v", err)
        } else {
            fmt.Println("谷歌回复提交成功")
        }
    }
}
```

## API 参考

### 类型

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

#### `ReviewItem`（响应结构）

| 字段 | 类型 | 说明 |
|------|------|------|
| `Platform` | `uint8` | `1`=Apple, `2`=Google |
| `ReviewId` | `string` | 评论唯一 ID |
| `ReviewTitle` | `string` | 评论标题 |
| `OriginalContent` | `string` | 原始评论内容 |
| `TranslatedContent` | `string` | 翻译内容（仅 Google Play 可能有） |
| `ReviewNickname` | `string` | 评论者昵称 |
| `ReviewRating` | `string` | 星级（"1"~"5"） |
| `ReviewLanguage` | `string` | 语言/地区代码 |
| `ReviewExtra` | `map[string]any` | 附加信息（谷歌含 device/app_version 等） |
| `CreatedAt` | `string` | 评论时间（RFC3339 格式） |

### 方法

| 方法 | 说明 |
|------|------|
| `New(ctx, apple, google) (*Client, error)` | 创建 SDK Client |
| `FetchAppleReviews(ctx, limit, since) ([]ReviewItem, error)` | 获取苹果评论 |
| `FetchGoogleReviews(ctx, limit, since) ([]ReviewItem, error)` | 获取谷歌评论 |
| `ReplyAppleReview(ctx, reviewID, content) error` | 提交苹果回复 |
| `ReplyGoogleReview(ctx, reviewID, content) error` | 提交谷歌回复 |

> `limit=0` 表示按时间窗口（苹果 48h / 谷歌 2h）拉取全部新评论；`limit>0` 表示最多拉取指定数量。

## 平台常量

```go
appstore.PlatformApple  // uint8 = 1
appstore.PlatformGoogle // uint8 = 2
```

## License

MIT License - see [LICENSE](LICENSE) for details

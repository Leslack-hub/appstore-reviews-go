# Getting Started

本指南将帮助您快速开始使用 appstore-reviews-go SDK。

## 前置要求

- Go 1.21 或更高版本
- Apple App Store Connect API 凭证（如果使用 Apple 功能）
- Google Play Developer API 凭证（如果使用 Google 功能）

## 安装

```bash
go get github.com/Leslack-hub/appstore-reviews-go
```

## 准备凭证

### Apple App Store Connect

1. 登录 [App Store Connect](https://appstoreconnect.apple.com/)
2. 导航到 "用户和访问" → "密钥" → "App Store Connect API"
3. 点击 "+" 创建新密钥
4. 选择适当的访问权限（需要 "Customer Reviews" 权限）
5. 下载 `.p8` 私钥文件（只能下载一次）
6. 记录以下信息：
   - Issuer ID（在密钥页面顶部）
   - Key ID（密钥列表中）
   - App ID（在 App Store Connect 中您的应用信息页面）

### Google Play Developer API

1. 登录 [Google Cloud Console](https://console.cloud.google.com/)
2. 创建新项目或选择现有项目
3. 启用 "Google Play Android Developer API"
4. 创建服务账号：
   - 导航到 "IAM & Admin" → "Service Accounts"
   - 点击 "Create Service Account"
   - 填写名称和描述
   - 点击 "Create and Continue"
5. 创建密钥：
   - 在服务账号列表中，点击刚创建的账号
   - 转到 "Keys" 标签
   - 点击 "Add Key" → "Create new key"
   - 选择 JSON 格式
   - 下载 JSON 文件
6. 在 [Google Play Console](https://play.google.com/console/) 中授权：
   - 导航到 "设置" → "API 访问"
   - 找到刚创建的服务账号
   - 点击 "授予访问权限"
   - 选择适当的权限（需要 "查看应用信息和下载批量报告" 和 "回复评价"）

## 第一个程序

创建 `main.go`：

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
    // 读取凭证文件
    certPEM, err := os.ReadFile("AuthKey_XXXXXXXXXX.p8")
    if err != nil {
        log.Fatal(err)
    }

    // 创建客户端
    client, err := appstore.NewAppleOnly(appstore.AppleConfig{
        IssuerID: "your-issuer-id",
        KeyID:    "your-key-id",
        Cert:     string(certPEM),
        AppID:    "your-app-id",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 获取评论
    ctx := context.Background()
    reviews, err := client.FetchAppleReviews(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    // 打印评论
    fmt.Printf("获取到 %d 条评论\n", len(reviews))
    for _, r := range reviews {
        fmt.Printf("[%s] %s ★%s: %s\n", 
            r.ReviewLanguage, 
            r.ReviewNickname, 
            r.ReviewRating, 
            r.OriginalContent)
    }
}
```

运行：

```bash
go run main.go
```

## 常见使用场景

### 场景 1: 获取最近的低分评论

```go
reviews, _ := client.FetchAppleReviews(ctx, nil)

for _, r := range reviews {
    rating, _ := strconv.Atoi(r.ReviewRating)
    if rating <= 2 {
        fmt.Printf("低分评论: %s\n", r.OriginalContent)
        // 可以在这里触发告警或自动回复
    }
}
```

### 场景 2: 自动回复评论

```go
reviews, _ := client.FetchAppleReviews(ctx, nil)

for _, r := range reviews {
    // 检查是否需要回复
    if shouldReply(r) {
        content := generateReply(r)
        err := client.ReplyAppleReview(ctx, r.ReviewId, content)
        if err != nil {
            log.Printf("回复失败: %v", err)
        }
    }
}
```

### 场景 3: 定期拉取评论

```go
import "time"

ticker := time.NewTicker(1 * time.Hour)
defer ticker.Stop()

for range ticker.C {
    reviews, err := client.FetchAppleReviews(ctx, &appstore.FetchAppleOptions{
        Since: 1 * time.Hour, // 只获取最近 1 小时的评论
    })
    if err != nil {
        log.Printf("获取评论失败: %v", err)
        continue
    }
    
    processReviews(reviews)
}
```

### 场景 4: 双平台支持

```go
credJSON, _ := os.ReadFile("google-service-account.json")

client, _ := appstore.New(ctx,
    appstore.AppleConfig{
        IssuerID: "apple-issuer-id",
        KeyID:    "apple-key-id",
        Cert:     string(certPEM),
        AppID:    "apple-app-id",
    },
    appstore.GoogleConfig{
        CredentialsJSON: credJSON,
        PackageName:     "com.example.app",
    },
)

// 同时获取两个平台的评论
appleReviews, _ := client.FetchAppleReviews(ctx, nil)
googleReviews, _ := client.FetchGoogleReviews(ctx, nil)

allReviews := append(appleReviews, googleReviews...)
```

## 最佳实践

### 1. 使用环境变量管理凭证

```go
import "os"

client, _ := appstore.NewAppleOnly(appstore.AppleConfig{
    IssuerID: os.Getenv("APPLE_ISSUER_ID"),
    KeyID:    os.Getenv("APPLE_KEY_ID"),
    Cert:     os.Getenv("APPLE_CERT"),
    AppID:    os.Getenv("APPLE_APP_ID"),
})
```

### 2. 使用 Context 控制超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

reviews, err := client.FetchAppleReviews(ctx, nil)
```

### 3. 实现错误重试

```go
func fetchWithRetry(client *appstore.Client, maxRetries int) ([]appstore.ReviewItem, error) {
    var reviews []appstore.ReviewItem
    var err error
    
    for i := 0; i < maxRetries; i++ {
        reviews, err = client.FetchAppleReviews(context.Background(), nil)
        if err == nil {
            return reviews, nil
        }
        
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return nil, err
}
```

### 4. 记录日志

```go
import "log"

reviews, err := client.FetchAppleReviews(ctx, nil)
if err != nil {
    log.Printf("获取评论失败: %v", err)
    return
}

log.Printf("成功获取 %d 条评论", len(reviews))
```

## 下一步

- 查看 [examples/](../examples/) 目录了解更多示例
- 阅读 [API 文档](https://pkg.go.dev/github.com/Leslack-hub/appstore-reviews-go)
- 查看 [FAQ](FAQ.md) 解决常见问题
- 了解 [架构设计](ARCHITECTURE.md)

## 获取帮助

如有问题：
- 查看 [FAQ](FAQ.md)
- 创建 [GitHub Issue](https://github.com/Leslack-hub/appstore-reviews-go/issues)
- 参与 [Discussions](https://github.com/Leslack-hub/appstore-reviews-go/discussions)

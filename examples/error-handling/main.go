package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	appstore "github.com/Leslack-hub/appstore-reviews-go"
)

func main() {
	// 示例 1: 配置验证错误
	handleConfigError()

	// 示例 2: 平台未配置错误
	handlePlatformNotConfigured()

	// 示例 3: API 调用错误
	handleAPIError()
}

func handleConfigError() {
	fmt.Println("=== 配置验证示例 ===")

	// 无效配置
	_, err := appstore.NewAppleOnly(appstore.AppleConfig{
		IssuerID: "issuer-id",
		// 缺少其他必需字段
	})

	if err != nil {
		if errors.Is(err, appstore.ErrInvalidConfig) {
			fmt.Println("配置无效，请检查所有必需字段")
		}
	}
}

func handlePlatformNotConfigured() {
	fmt.Println("\n=== 平台未配置示例 ===")

	// 创建仅支持 Apple 的客户端
	client, _ := appstore.NewAppleOnly(appstore.AppleConfig{
		IssuerID: "issuer-id",
		KeyID:    "key-id",
		Cert:     "cert-content",
		AppID:    "123456",
	})

	ctx := context.Background()

	// 尝试调用 Google 方法
	_, err := client.FetchGoogleReviews(ctx, nil)
	if err != nil {
		if errors.Is(err, appstore.ErrGoogleNotConfigured) {
			fmt.Println("Google 平台未配置")
		}
	}
}

func handleAPIError() {
	fmt.Println("\n=== API 调用错误处理示例 ===")

	certPEM, err := os.ReadFile("AuthKey_XXXXXXXXXX.p8")
	if err != nil {
		log.Printf("无法读取凭证文件: %v", err)
		return
	}

	client, err := appstore.NewAppleOnly(appstore.AppleConfig{
		IssuerID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		KeyID:    "XXXXXXXXXX",
		Cert:     string(certPEM),
		AppID:    "6737681541",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 获取评论，处理可能的错误
	reviews, err := client.FetchAppleReviews(ctx, nil)
	if err != nil {
		// 根据错误类型采取不同的处理策略
		switch {
		case errors.Is(err, context.Canceled):
			fmt.Println("操作被取消")
		case errors.Is(err, context.DeadlineExceeded):
			fmt.Println("操作超时")
		default:
			fmt.Printf("API 调用失败: %v\n", err)
			// 可以在这里实现重试逻辑
		}
		return
	}

	fmt.Printf("成功获取 %d 条评论\n", len(reviews))
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	appstore "github.com/Leslack-hub/appstore-reviews-go"
)

func main() {
	// 示例 1: 仅使用 Apple
	appleOnly()

	// 示例 2: 仅使用 Google
	googleOnly()
}

func appleOnly() {
	certPEM, err := os.ReadFile("AuthKey_XXXXXXXXXX.p8")
	if err != nil {
		log.Printf("跳过 Apple 示例: %v", err)
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

	// 检查平台状态
	fmt.Printf("Apple enabled: %v\n", client.IsAppleEnabled())
	fmt.Printf("Google enabled: %v\n", client.IsGoogleEnabled())

	// 获取评论
	reviews, err := client.FetchAppleReviews(ctx, nil)
	if err != nil {
		log.Printf("获取评论失败: %v", err)
		return
	}

	fmt.Printf("获取到 %d 条 Apple 评论\n", len(reviews))
}

func googleOnly() {
	credJSON, err := os.ReadFile("google-service-account.json")
	if err != nil {
		log.Printf("跳过 Google 示例: %v", err)
		return
	}

	client, err := appstore.NewGoogleOnly(context.Background(), appstore.GoogleConfig{
		CredentialsJSON: credJSON,
		PackageName:     "com.example.myapp",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 检查平台状态
	fmt.Printf("Apple enabled: %v\n", client.IsAppleEnabled())
	fmt.Printf("Google enabled: %v\n", client.IsGoogleEnabled())

	// 获取评论
	reviews, err := client.FetchGoogleReviews(ctx, nil)
	if err != nil {
		log.Printf("获取评论失败: %v", err)
		return
	}

	fmt.Printf("获取到 %d 条 Google 评论\n", len(reviews))
}

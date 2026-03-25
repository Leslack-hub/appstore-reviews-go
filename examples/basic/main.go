package main

import (
	"context"
	"fmt"
	"log"
	"os"

	appstore "github.com/Leslack-hub/appstore-reviews-go"
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

	// 获取苹果评论（最近 48 小时）
	appleReviews, err := client.FetchAppleReviews(ctx, nil)
	if err != nil {
		log.Printf("获取苹果评论失败: %v", err)
	} else {
		fmt.Printf("苹果评论数: %d\n", len(appleReviews))
		for _, r := range appleReviews {
			fmt.Printf("  [%s] %s ★%s: %s\n", r.ReviewLanguage, r.ReviewNickname, r.ReviewRating, r.OriginalContent)
		}
	}

	// 获取谷歌评论（最近 2 小时）
	googleReviews, err := client.FetchGoogleReviews(ctx, nil)
	if err != nil {
		log.Printf("获取谷歌评论失败: %v", err)
	} else {
		fmt.Printf("谷歌评论数: %d\n", len(googleReviews))
	}

	// 提交苹果回复
	if len(appleReviews) > 0 {
		r := appleReviews[0]
		if err := client.ReplyAppleReview(ctx, r.ReviewId, "感谢您的反馈！"); err != nil {
			log.Printf("回复苹果评论失败: %v", err)
		} else {
			fmt.Println("苹果回复提交成功")
		}
	}

	// 提交谷歌回复
	if len(googleReviews) > 0 {
		r := googleReviews[0]
		if err := client.ReplyGoogleReview(ctx, r.ReviewId, "感谢您的反馈！"); err != nil {
			log.Printf("回复谷歌评论失败: %v", err)
		} else {
			fmt.Println("谷歌回复提交成功")
		}
	}
}

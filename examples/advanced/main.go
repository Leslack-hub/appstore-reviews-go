package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	appstore "github.com/Leslack-hub/appstore-reviews-go"
)

func main() {
	certPEM, _ := os.ReadFile("AuthKey_XXXXXXXXXX.p8")

	// 创建仅支持 Apple 的客户端
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

	// 使用高级选项获取评论
	reviews, err := client.FetchAppleReviews(ctx, &appstore.FetchAppleOptions{
		Limit:   10,                 // 最多获取 10 条
		Since:   7 * 24 * time.Hour, // 最近 7 天
		PerPage: 50,                 // 每页 50 条
		Sort:    "-createdDate",     // 按创建时间倒序
		OnPage: func(items []appstore.ReviewItem) bool {
			fmt.Printf("收到 %d 条评论\n", len(items))
			return true // 继续获取下一页
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("总共获取 %d 条评论\n", len(reviews))

	// 只处理低分评论
	for _, r := range reviews {
		if r.ReviewRating <= "2" {
			fmt.Printf("低分评论: ★%s - %s\n", r.ReviewRating, r.OriginalContent)
			// 可以在这里自动回复
		}
	}
}

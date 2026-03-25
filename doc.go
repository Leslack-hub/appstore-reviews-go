// Package appstore 提供从 Apple App Store 和 Google Play 获取用户评论及提交开发者回复的 SDK。
//
// 该包封装了 Apple App Store Connect API 和 Google Play Developer API 的调用，
// 提供统一的接口来管理应用商店评论。
//
// # 基本用法
//
// 创建客户端并获取评论：
//
//	client, err := appstore.New(ctx,
//		appstore.AppleConfig{
//			IssuerID: "your-issuer-id",
//			KeyID:    "your-key-id",
//			Cert:     certPEM,
//			AppID:    "your-app-id",
//		},
//		appstore.GoogleConfig{
//			CredentialsJSON: credJSON,
//			PackageName:     "com.example.app",
//		},
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 获取最近 48 小时的 Apple 评论
//	reviews, err := client.FetchAppleReviews(ctx, 0, 0)
//
//	// 回复评论
//	err = client.ReplyAppleReview(ctx, reviewID, "感谢您的反馈！")
//
// # 平台支持
//
// - Apple App Store: 使用 App Store Connect API，需要 .p8 私钥文件
// - Google Play: 使用 Google Play Developer API，需要服务账号 JSON 凭证
//
// # 时间窗口
//
// 当 limit 参数为 0 时，SDK 会按默认时间窗口获取评论：
//   - Apple: 最近 48 小时
//   - Google: 最近 2 小时
package appstore

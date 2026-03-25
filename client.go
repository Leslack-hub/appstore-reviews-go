package appstore

import (
	"context"
	"time"
)

// AppleConfig 包含访问 Apple App Store Connect API 所需的配置信息。
type AppleConfig struct {
	// IssuerID 是 App Store Connect API 的 Issuer ID
	IssuerID string
	// KeyID 是 API Key ID
	KeyID string
	// Cert 是 ECDSA 私钥的 PEM 格式内容（.p8 文件内容）
	Cert string
	// AppID 是苹果应用的数字 ID
	AppID string
}

// GoogleConfig 包含访问 Google Play Developer API 所需的配置信息。
type GoogleConfig struct {
	// CredentialsJSON 是服务账号 JSON 文件的内容
	CredentialsJSON []byte
	// PackageName 是 Android 应用的包名
	PackageName string
}

// Client 是应用商店评论管理的客户端，支持 Apple App Store 和 Google Play。
type Client struct {
	apple  AppleConfig
	google *googleClient
}

// New 创建一个新的 Client 实例。
//
// 参数：
//   - ctx: 上下文，用于控制初始化过程
//   - apple: Apple App Store Connect API 配置
//   - googleCfg: Google Play Developer API 配置
//
// 如果 googleCfg.CredentialsJSON 为空或不是有效的 JSON，Google Play 功能将不可用。
// 返回的 Client 可以安全地调用所有方法，未配置的平台方法将返回 nil 或错误。
func New(ctx context.Context, apple AppleConfig, googleCfg GoogleConfig) (*Client, error) {
	c := &Client{
		apple: apple,
	}
	if len(googleCfg.CredentialsJSON) > 0 && googleCfg.CredentialsJSON[0] == '{' {
		gc, err := newGoogleClient(ctx, googleCfg)
		if err != nil {
			return nil, err
		}
		c.google = gc
	}
	return c, nil
}

// FetchAppleReviews 获取 Apple App Store 的用户评论。
//
// 参数：
//   - ctx: 上下文，用于控制请求过程
//   - limit: 最多获取的评论数量。如果为 0，则获取指定时间窗口内的所有评论
//   - since: 时间窗口，获取该时间段内的评论。如果为 0，默认使用 48 小时
//
// 返回评论列表按创建时间倒序排列（最新的在前）。
// 如果 Apple 配置未设置（IssuerID 为空），返回 nil, nil。
func (c *Client) FetchAppleReviews(ctx context.Context, limit int, since time.Duration) ([]ReviewItem, error) {
	if c.apple.IssuerID == "" {
		return nil, nil
	}
	return fetchAppleReviews(ctx, c.apple, limit, since)
}

// FetchGoogleReviews 获取 Google Play 的用户评论。
//
// 参数：
//   - ctx: 上下文，用于控制请求过程
//   - limit: 最多获取的评论数量。如果为 0，则获取指定时间窗口内的所有评论
//   - since: 时间窗口，获取该时间段内的评论。如果为 0，默认使用 2 小时
//
// 返回评论列表按最后修改时间排序。
// 如果 Google 配置未设置，返回 nil, nil。
func (c *Client) FetchGoogleReviews(ctx context.Context, limit int, since time.Duration) ([]ReviewItem, error) {
	if c.google == nil {
		return nil, nil
	}
	return c.google.fetchReviews(ctx, limit, since)
}

// ReplyAppleReview 向 Apple App Store 的评论提交开发者回复。
//
// 参数：
//   - ctx: 上下文，用于控制请求过程
//   - reviewID: 评论的唯一标识符
//   - content: 回复内容
//
// 注意：每条评论只能回复一次，重复回复会返回错误。
func (c *Client) ReplyAppleReview(ctx context.Context, reviewID, content string) error {
	return replyAppleReview(ctx, c.apple, reviewID, content)
}

// ReplyGoogleReview 向 Google Play 的评论提交开发者回复。
//
// 参数：
//   - ctx: 上下文，用于控制请求过程
//   - reviewID: 评论的唯一标识符
//   - content: 回复内容
//
// 如果 Google 配置未设置，返回 context.Canceled 错误。
// Google Play 允许多次回复同一条评论，新回复会覆盖旧回复。
func (c *Client) ReplyGoogleReview(ctx context.Context, reviewID, content string) error {
	if c.google == nil {
		return context.Canceled
	}
	return c.google.replyReview(ctx, reviewID, content)
}

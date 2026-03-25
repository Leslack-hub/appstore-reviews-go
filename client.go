package appstore

import (
	"context"
	"net/url"
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

// FetchAppleOptions 定义了拉取苹果评论的高级配置。
type FetchAppleOptions struct {
	// 期望获取的最大评论数量。0 表示不限制
	Limit int
	// 期望获取的时间窗口。0 表示不限制时间
	Since time.Duration
	// 每页拉取数量 (1-200)，默认 200
	PerPage int
	// 排序方式，默认 "-createdDate"
	Sort string
	// 额外的查询参数（例如 url.Values{"filter[rating]": {"1"}} 过滤一星评论）
	QueryParams url.Values
	// 接收到一页数据后的回调函数。如果返回 false，将中止拉取
	OnPage func(items []ReviewItem) bool
}

// FetchGoogleOptions 定义了拉取谷歌评论的高级配置。
type FetchGoogleOptions struct {
	// 期望获取的最大评论数量。0 表示不限制
	Limit int
	// 期望获取的时间窗口。0 表示不限制时间
	Since time.Duration
	// 翻译语言参数
	TranslationLanguage string
	// 接收到一页数据后的回调函数。如果返回 false，将中止拉取
	OnPage func(items []ReviewItem) bool
}

// FetchAppleReviews 获取 Apple App Store 的用户评论。
//
// 参数：
//   - ctx: 上下文，用于控制请求过程
//   - opts: 高级查询配置，提供极大的灵活性
//
// 如果 Apple 配置未设置（IssuerID 为空），返回 nil, nil。
func (c *Client) FetchAppleReviews(ctx context.Context, opts *FetchAppleOptions) ([]ReviewItem, error) {
	if c.apple.IssuerID == "" {
		return nil, nil
	}
	if opts == nil {
		opts = &FetchAppleOptions{Since: 48 * time.Hour}
	}
	return fetchAppleReviews(ctx, c.apple, opts)
}

// FetchGoogleReviews 获取 Google Play 的用户评论。
//
// 参数：
//   - ctx: 上下文，用于控制请求过程
//   - opts: 高级查询配置，提供极大的灵活性
//
// 如果 Google 配置未设置，返回 nil, nil。
func (c *Client) FetchGoogleReviews(ctx context.Context, opts *FetchGoogleOptions) ([]ReviewItem, error) {
	if c.google == nil {
		return nil, nil
	}
	if opts == nil {
		opts = &FetchGoogleOptions{Since: 2 * time.Hour}
	}
	return c.google.fetchReviews(ctx, opts)
}

// ReplyAppleReview 向 Apple App Store 的评论提交开发者回复。
func (c *Client) ReplyAppleReview(ctx context.Context, reviewID, content string) error {
	return replyAppleReview(ctx, c.apple, reviewID, content)
}

// ReplyGoogleReview 向 Google Play 的评论提交开发者回复。
func (c *Client) ReplyGoogleReview(ctx context.Context, reviewID, content string) error {
	if c.google == nil {
		return context.Canceled
	}
	return c.google.replyReview(ctx, reviewID, content)
}

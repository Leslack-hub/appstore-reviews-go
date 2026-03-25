package appstore

const (
	// PlatformApple 表示 iOS App Store 平台
	PlatformApple uint8 = 1
	// PlatformGoogle 表示 Google Play 平台
	PlatformGoogle uint8 = 2
)

// ReviewItem 统一的评论条目，供 SDK 返回及 HTTP 传输使用。
// 字段与苹果/谷歌实际提供的数据严格对齐，CreatedAt 使用 RFC3339 字符串。
type ReviewItem struct {
	// Platform 平台: 1=Apple App Store, 2=Google Play
	Platform uint8 `json:"platform"`
	// ReviewId 评论唯一标识
	ReviewId string `json:"review_id"`
	// ReviewTitle 评论标题（苹果有，谷歌通常为空）
	ReviewTitle string `json:"review_title"`
	// OriginalContent 评论原始内容
	OriginalContent string `json:"original_content"`
	// TranslatedContent 谷歌翻译内容（仅 Google Play 在原文非英文时提供，苹果为空）
	TranslatedContent string `json:"translated_content"`
	// ReviewNickname 评论者昵称
	ReviewNickname string `json:"review_nickname"`
	// ReviewRating 评论星级字符串（"1"~"5"）
	ReviewRating string `json:"review_rating"`
	// ReviewLanguage 评论语言/地区（苹果为 territory 代码，谷歌为 BCP-47 标签）
	ReviewLanguage string `json:"review_language"`
	// ReviewExtra 谷歌附加信息（device/device_meta/app_version/version_name/os_version），苹果为空 map
	ReviewExtra map[string]any `json:"review_extra,omitempty"`
	// CreatedAt 评论时间，RFC3339 格式
	CreatedAt string `json:"created_at"`
}

// FetchReviewDataResp 定义对外提供带有 Prompt 的评论列表数据
type FetchReviewDataResp struct {
	List          []ReviewItem `json:"list"`
	CommentPrompt string       `json:"comment_prompt"`
	PolishPrompt  string       `json:"polish_prompt"`
}

// FetchReviewResponse 约定各个微服务业务方返回的标准含 Prompt 外壳的响应格式
type FetchReviewResponse struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data FetchReviewDataResp `json:"data"`
}

// SubmitReplyRequest 约定 backend 向各项目组回传 AI 回复的请求体
type SubmitReplyRequest struct {
	ReviewId     string `json:"review_id"`
	ReplyContent string `json:"reply_content"`
}

package appstore

import "errors"

var (
	// ErrAppleNotConfigured 表示 Apple 配置未设置
	ErrAppleNotConfigured = errors.New("apple configuration not set")
	// ErrGoogleNotConfigured 表示 Google 配置未设置
	ErrGoogleNotConfigured = errors.New("google configuration not set")
	// ErrInvalidConfig 表示配置参数无效
	ErrInvalidConfig = errors.New("invalid configuration")
	// ErrInvalidReviewID 表示评论 ID 无效
	ErrInvalidReviewID = errors.New("invalid review ID")
	// ErrEmptyReplyContent 表示回复内容为空
	ErrEmptyReplyContent = errors.New("reply content cannot be empty")
)

package appstore

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type googleClient struct {
	svc         *androidpublisher.Service
	packageName string
}

func newGoogleClient(ctx context.Context, cfg GoogleConfig) (*googleClient, error) {
	jwtCfg, err := google.JWTConfigFromJSON(cfg.CredentialsJSON, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, fmt.Errorf("google: parse credentials: %w", err)
	}
	svc, err := androidpublisher.NewService(ctx, option.WithHTTPClient(jwtCfg.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("google: create service: %w", err)
	}
	return &googleClient{svc: svc, packageName: cfg.PackageName}, nil
}

func (g *googleClient) fetchReviews(ctx context.Context, opts *FetchGoogleOptions) ([]ReviewItem, error) {
	var cutoff time.Time
	if opts.Since > 0 {
		cutoff = time.Now().Add(-opts.Since)
	}
	var results []ReviewItem
	var nextPageToken string

Label1:
	for {
		select {
		case <-ctx.Done():
			break Label1
		default:
		}
		call := g.svc.Reviews.List(g.packageName)
		if nextPageToken != "" {
			call.Token(nextPageToken)
		}
		if opts.TranslationLanguage != "" {
			call.TranslationLanguage(opts.TranslationLanguage)
		}

		resp, err := call.Do()
		if err != nil {
			return results, fmt.Errorf("google: fetch reviews: %w", err)
		}

		var pageItems []ReviewItem
		shouldBreak := false
		for _, review := range resp.Reviews {
			if opts.Limit > 0 && len(results)+len(pageItems) >= opts.Limit {
				shouldBreak = true
				break
			}
			if len(review.Comments) == 0 || review.Comments[0].UserComment == nil {
				continue
			}

			userComment := review.Comments[0].UserComment
			var reviewTime time.Time
			if userComment.LastModified != nil {
				reviewTime = time.Unix(userComment.LastModified.Seconds, userComment.LastModified.Nanos)
			}
			if !cutoff.IsZero() && !reviewTime.IsZero() && reviewTime.Before(cutoff) {
				shouldBreak = true
				break
			}

			originalContent := userComment.Text
			translatedContent := ""
			if userComment.OriginalText != "" {
				originalContent = userComment.OriginalText
				translatedContent = userComment.Text
			}

			reviewTitle := ""
			parts := strings.SplitN(originalContent, "\t", 2)
			if len(parts) == 2 {
				reviewTitle = strings.TrimSpace(parts[0])
				originalContent = strings.TrimSpace(parts[1])
			}

			pageItems = append(pageItems, ReviewItem{
				Platform:          PlatformGoogle,
				ReviewId:          review.ReviewId,
				ReviewTitle:       reviewTitle,
				OriginalContent:   originalContent,
				TranslatedContent: translatedContent,
				ReviewNickname:    review.AuthorName,
				ReviewRating:      strconv.Itoa(int(userComment.StarRating)),
				ReviewLanguage:    userComment.ReviewerLanguage,
				ReviewExtra: map[string]any{
					"device":       userComment.Device,
					"device_meta":  userComment.DeviceMetadata,
					"app_version":  userComment.AppVersionCode,
					"version_name": userComment.AppVersionName,
					"os_version":   userComment.AndroidOsVersion,
				},
				CreatedAt: reviewTime.Format(time.RFC3339),
			})
		}
		results = append(results, pageItems...)
		if opts.OnPage != nil && len(pageItems) > 0 {
			if !opts.OnPage(pageItems) {
				break
			}
		}
		if shouldBreak || resp.TokenPagination == nil || resp.TokenPagination.NextPageToken == "" {
			break
		}
		nextPageToken = resp.TokenPagination.NextPageToken
	}

	return results, nil
}

func (g *googleClient) replyReview(_ context.Context, reviewID, content string) error {
	_, err := g.svc.Reviews.Reply(g.packageName, reviewID, &androidpublisher.ReviewsReplyRequest{
		ReplyText: content,
	}).Do()
	if err != nil {
		return fmt.Errorf("google reply: api call: %w", err)
	}
	return nil
}

package appstore

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jwtpkg "github.com/golang-jwt/jwt/v4"
)

const appleConnectBaseURL = "https://api.appstoreconnect.apple.com"

func generateAppleToken(cfg AppleConfig) (string, error) {
	block, _ := pem.Decode([]byte(cfg.Cert))
	if block == nil {
		return "", fmt.Errorf("apple: failed to decode PEM block")
	}
	keyIface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("apple: parse private key: %w", err)
	}
	ecKey, ok := keyIface.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("apple: not an ECDSA key")
	}

	now := time.Now()
	tok := &jwtpkg.Token{
		Header: map[string]interface{}{"alg": "ES256", "kid": cfg.KeyID, "typ": "JWT"},
		Claims: jwtpkg.MapClaims{
			"iss": cfg.IssuerID,
			"iat": now.Unix(),
			"exp": now.Add(20 * time.Minute).Unix(),
			"aud": "appstoreconnect-v1",
		},
		Method: jwtpkg.SigningMethodES256,
	}
	return tok.SignedString(ecKey)
}

type appleReviewAttributes struct {
	Rating           int    `json:"rating"`
	Title            string `json:"title"`
	Body             string `json:"body"`
	ReviewerNickname string `json:"reviewerNickname"`
	CreatedDate      string `json:"createdDate"`
	Territory        string `json:"territory"`
}

type appleReviewData struct {
	ID         string                 `json:"id"`
	Attributes *appleReviewAttributes `json:"attributes"`
}

type appleReviewLinks struct {
	Next string `json:"next"`
}

type appleReviewResponse struct {
	Data  []appleReviewData `json:"data"`
	Links appleReviewLinks  `json:"links"`
}

func fetchAppleReviews(ctx context.Context, cfg AppleConfig, opts *FetchAppleOptions) ([]ReviewItem, error) {
	token, err := generateAppleToken(cfg)
	if err != nil {
		return nil, fmt.Errorf("apple: generate token: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// 构造查询参数
	q := url.Values{}
	if opts.Sort != "" {
		q.Set("sort", opts.Sort)
	} else {
		q.Set("sort", "-createdDate")
	}
	if opts.PerPage > 0 {
		q.Set("limit", strconv.Itoa(opts.PerPage))
	} else {
		q.Set("limit", "200")
	}
	for k, v := range opts.QueryParams {
		for _, val := range v {
			q.Add(k, val)
		}
	}

	baseURL := fmt.Sprintf("%s/v1/apps/%s/customerReviews?%s", appleConnectBaseURL, cfg.AppID, q.Encode())
	nextURL := baseURL

	var cutoff time.Time
	if opts.Since > 0 {
		cutoff = time.Now().Add(-opts.Since)
	}
	var results []ReviewItem

	for nextURL != "" {
		var req *http.Request
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, nextURL, nil)
		if err != nil {
			return results, fmt.Errorf("apple: build request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		var resp *http.Response
		resp, err = client.Do(req)
		if err != nil {
			return results, fmt.Errorf("apple: http request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return results, fmt.Errorf("apple: unexpected status %d: %s", resp.StatusCode, body)
		}

		var result appleReviewResponse
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			_ = resp.Body.Close()
			return results, fmt.Errorf("apple: decode response: %w", err)
		}
		_ = resp.Body.Close()

		var pageItems []ReviewItem
		shouldBreak := false
		for _, review := range result.Data {
			a := review.Attributes
			if opts.Limit > 0 && len(results)+len(pageItems) >= opts.Limit {
				shouldBreak = true
				break
			}

			if !cutoff.IsZero() {
				var createdAt time.Time
				if createdAt, err = time.Parse(time.RFC3339, a.CreatedDate); err == nil && createdAt.Before(cutoff) {
					shouldBreak = true
					break
				}
			}

			pageItems = append(pageItems, ReviewItem{
				Platform:        PlatformApple,
				ReviewId:        review.ID,
				ReviewTitle:     a.Title,
				OriginalContent: a.Body,
				ReviewNickname:  a.ReviewerNickname,
				ReviewRating:    strconv.Itoa(a.Rating),
				ReviewLanguage:  a.Territory,
				CreatedAt:       a.CreatedDate,
			})
		}

		results = append(results, pageItems...)
		if opts.OnPage != nil && len(pageItems) > 0 {
			if !opts.OnPage(pageItems) {
				break
			}
		}
		if shouldBreak {
			break
		}
		nextURL = result.Links.Next
	}

	return results, nil
}

func replyAppleReview(ctx context.Context, cfg AppleConfig, reviewID, content string) error {
	token, err := generateAppleToken(cfg)
	if err != nil {
		return fmt.Errorf("apple reply: generate token: %w", err)
	}

	type relData struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	type rel struct {
		Review struct {
			Data relData `json:"data"`
		} `json:"review"`
	}
	type attrs struct {
		ResponseBody string `json:"responseBody"`
	}
	type dataBody struct {
		Type          string `json:"type"`
		Attributes    attrs  `json:"attributes"`
		Relationships rel    `json:"relationships"`
	}
	type reqBody struct {
		Data dataBody `json:"data"`
	}

	payload, err := json.Marshal(reqBody{
		Data: dataBody{
			Type:       "customerReviewResponses",
			Attributes: attrs{ResponseBody: content},
			Relationships: rel{
				Review: struct {
					Data relData `json:"data"`
				}{Data: relData{ID: reviewID, Type: "customerReviews"}},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("apple reply: marshal: %w", err)
	}

	url := appleConnectBaseURL + "/v1/customerReviewResponses"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("apple reply: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("apple reply: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("apple reply: unexpected status %d: %s", resp.StatusCode, body)
	}

	return nil
}

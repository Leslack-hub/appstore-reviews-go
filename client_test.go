package appstore

import (
	"context"
	"testing"
)

func TestAppleConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  AppleConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: AppleConfig{
				IssuerID: "issuer-id",
				KeyID:    "key-id",
				Cert:     "cert-content",
				AppID:    "123456",
			},
			wantErr: false,
		},
		{
			name: "missing issuer id",
			config: AppleConfig{
				KeyID: "key-id",
				Cert:  "cert-content",
				AppID: "123456",
			},
			wantErr: true,
		},
		{
			name: "missing key id",
			config: AppleConfig{
				IssuerID: "issuer-id",
				Cert:     "cert-content",
				AppID:    "123456",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  AppleConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AppleConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoogleConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  GoogleConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: GoogleConfig{
				CredentialsJSON: []byte(`{"type":"service_account"}`),
				PackageName:     "com.example.app",
			},
			wantErr: false,
		},
		{
			name: "missing credentials",
			config: GoogleConfig{
				PackageName: "com.example.app",
			},
			wantErr: true,
		},
		{
			name: "missing package name",
			config: GoogleConfig{
				CredentialsJSON: []byte(`{"type":"service_account"}`),
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  GoogleConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAppleOnly(t *testing.T) {
	tests := []struct {
		name    string
		config  AppleConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: AppleConfig{
				IssuerID: "issuer-id",
				KeyID:    "key-id",
				Cert:     "cert-content",
				AppID:    "123456",
			},
			wantErr: false,
		},
		{
			name:    "invalid config",
			config:  AppleConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewAppleOnly(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAppleOnly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if client == nil {
					t.Error("NewAppleOnly() returned nil client")
				}
				if !client.IsAppleEnabled() {
					t.Error("Apple should be enabled")
				}
				if client.IsGoogleEnabled() {
					t.Error("Google should not be enabled")
				}
			}
		})
	}
}

func TestClient_FetchAppleReviews_NotConfigured(t *testing.T) {
	client := &Client{appleEnabled: false}
	ctx := context.Background()

	_, err := client.FetchAppleReviews(ctx, nil)
	if err != ErrAppleNotConfigured {
		t.Errorf("expected ErrAppleNotConfigured, got %v", err)
	}
}

func TestClient_FetchGoogleReviews_NotConfigured(t *testing.T) {
	client := &Client{google: nil}
	ctx := context.Background()

	_, err := client.FetchGoogleReviews(ctx, nil)
	if err != ErrGoogleNotConfigured {
		t.Errorf("expected ErrGoogleNotConfigured, got %v", err)
	}
}

func TestClient_ReplyAppleReview_NotConfigured(t *testing.T) {
	client := &Client{appleEnabled: false}
	ctx := context.Background()

	err := client.ReplyAppleReview(ctx, "review-id", "content")
	if err != ErrAppleNotConfigured {
		t.Errorf("expected ErrAppleNotConfigured, got %v", err)
	}
}

func TestClient_ReplyGoogleReview_NotConfigured(t *testing.T) {
	client := &Client{google: nil}
	ctx := context.Background()

	err := client.ReplyGoogleReview(ctx, "review-id", "content")
	if err != ErrGoogleNotConfigured {
		t.Errorf("expected ErrGoogleNotConfigured, got %v", err)
	}
}

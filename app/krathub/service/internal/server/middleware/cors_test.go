package middleware

import (
	"testing"
	"time"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

func TestCORSConfigFromConfig_Disabled(t *testing.T) {
	// Test when config is nil
	options := CORS(nil)
	if len(options.AllowedOrigins) != 0 {
		t.Errorf("Expected empty options when config is nil, got %v", options)
	}

	// Test when enable is false
	config := &conf.CORS{
		Enable: false,
	}
	options = CORS(config)
	if len(options.AllowedOrigins) != 0 {
		t.Errorf("Expected empty options when enable is false, got %v", options)
	}
}

func TestCORSConfigFromConfig_Enabled(t *testing.T) {
	config := &conf.CORS{
		Enable:           true,
		AllowedOrigins:   []string{"https://example.com", "https://test.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Custom-Header"},
		AllowCredentials: true,
		MaxAge: &durationpb.Duration{
			Seconds: 3600,
		},
	}

	options := CORS(config)

	if len(options.AllowedOrigins) != 2 {
		t.Errorf("Expected 2 allowed origins, got %d", len(options.AllowedOrigins))
	}

	if options.AllowedOrigins[0] != "https://example.com" {
		t.Errorf("Expected first allowed origin to be 'https://example.com', got %s", options.AllowedOrigins[0])
	}

	if !options.AllowCredentials {
		t.Error("Expected allow credentials to be true")
	}

	if options.MaxAge != time.Hour {
		t.Errorf("Expected max age to be 1h, got %v", options.MaxAge)
	}
}

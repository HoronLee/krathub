package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOptions_DefaultValues(t *testing.T) {
	options := DefaultOptions()

	if len(options.AllowedOrigins) != 1 || options.AllowedOrigins[0] != "*" {
		t.Errorf("Expected default allowed origins to be ['*'], got %v", options.AllowedOrigins)
	}

	if len(options.AllowedMethods) != 5 {
		t.Errorf("Expected 5 default allowed methods, got %d", len(options.AllowedMethods))
	}

	if !contains(options.AllowedMethods, "GET") || !contains(options.AllowedMethods, "POST") {
		t.Errorf("Expected GET and POST in default allowed methods, got %v", options.AllowedMethods)
	}

	if options.AllowCredentials {
		t.Error("Expected default allow credentials to be false")
	}

	if options.MaxAge != 24*time.Hour {
		t.Errorf("Expected default max age to be 24h, got %v", options.MaxAge)
	}
}

func TestMiddleware_Disabled(t *testing.T) {
	options := Options{} // Empty options means disabled
	corsMiddleware := Middleware(options)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Error("Expected no CORS headers when disabled")
	}
}

func TestMiddleware_SimpleRequest(t *testing.T) {
	options := Options{
		AllowedOrigins:   []string{"https://example.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           time.Hour,
	}

	corsMiddleware := Middleware(options)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.Header.Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected 'https://example.com' in Access-Control-Allow-Origin, got %s", res.Header.Get("Access-Control-Allow-Origin"))
	}

	if res.Header.Get("Access-Control-Allow-Methods") != "GET, POST" {
		t.Errorf("Expected 'GET, POST' in Access-Control-Allow-Methods, got %s", res.Header.Get("Access-Control-Allow-Methods"))
	}

	if res.Header.Get("Access-Control-Allow-Credentials") == "true" {
		t.Error("Expected no Access-Control-Allow-Credentials header")
	}
}

func TestMiddleware_PreflightRequest(t *testing.T) {
	options := Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           time.Hour,
	}

	corsMiddleware := Middleware(options)

	req := httptest.NewRequest("OPTIONS", "http://example.com", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 for preflight request, got %d", res.StatusCode)
	}

	if res.Header.Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected 'https://example.com' in Access-Control-Allow-Origin, got %s", res.Header.Get("Access-Control-Allow-Origin"))
	}

	if res.Header.Get("Access-Control-Max-Age") != "3600" {
		t.Errorf("Expected '3600' in Access-Control-Max-Age, got %s", res.Header.Get("Access-Control-Max-Age"))
	}
}

func TestMiddleware_OriginNotAllowed(t *testing.T) {
	options := Options{
		AllowedOrigins: []string{"https://allowed.com"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type"},
	}

	corsMiddleware := Middleware(options)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Origin", "https://notallowed.com")
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	res := w.Result()
	if res.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Error("Expected no Access-Control-Allow-Origin header when origin is not allowed")
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name          string
		origin        string
		allowedOrigin []string
		expected      bool
	}{
		{"wildcard", "https://example.com", []string{"*"}, true},
		{"exact match", "https://example.com", []string{"https://example.com"}, true},
		{"no match", "https://example.com", []string{"https://different.com"}, false},
		{"empty origin", "", []string{"*"}, false},
		{"wildcard subdomain", "https://api.example.com", []string{"*.example.com"}, true},
		{"wildcard subdomain no match", "https://api.baddomain.com", []string{"*.example.com"}, false},
		{"wildcard subdomain too many dots", "https://fake.api.example.com", []string{"*.example.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowedOrigin)
			if result != tt.expected {
				t.Errorf("Expected %v for origin %s in %v, got %v", tt.expected, tt.origin, tt.allowedOrigin, result)
			}
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

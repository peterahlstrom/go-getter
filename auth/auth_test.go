package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peterahlstrom/go-getter/config"
)

func TestApiKeyMiddleware(t *testing.T) {
	endpoints := map[string]config.Endpoint{
		"/public": {
			RequireAuth: false,
		},
		"/protected": {
			RequireAuth: true,
			ValidApiKeys: map[string]string{
				"abc123": "dev",
				"def456": "prod",
			},
		},
	}

	mw := ApiKeyMiddleWare(endpoints)

	tests := []struct {
		name       string
		url        string
		authHeader string
		wantStatus int
	}{
		{"invalid endpoint", "/nope", "", http.StatusNotFound},
		{"no auth", "/public", "", http.StatusOK},
		{"missing key", "/protected", "", http.StatusUnauthorized},
		{"invalid key format", "/protected", "abc123", http.StatusUnauthorized},
		{"invalid key", "/protected", "ApiKey nope", http.StatusUnauthorized},
		{"valid key", "/protected", "ApiKey abc123", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

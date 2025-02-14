package auth

import (
	"net/http"
	"strings"

	"github.com/peterahlstrom/go-getter/config"
)

type ApiKeyError struct {
	Message string
	HttpStatus int
}

func (e *ApiKeyError) Error() string {
	return e.Message
}

func ApiKeyMiddleWare(endpoints map[string]config.Endpoint) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			e, exists := endpoints[r.URL.Path]
			if !exists {
				http.Error(w, "Invalid endpoint", http.StatusNotFound)
				return
			}
		
			if !e.RequireAuth {
				next.ServeHTTP(w, r)
				return
			}
		
			reqKey := r.Header.Get("Authorization")
			if reqKey == "" {
				http.Error(w, "Missing API key", http.StatusUnauthorized)
				return
			}
		
			const prefix = "ApiKey "
			if !strings.HasPrefix(reqKey, prefix) {
				http.Error(w, "Invalid API key format", http.StatusUnauthorized)
				return
			}
		
			key := strings.TrimPrefix(reqKey, prefix)
			if _, valid := e.ValidApiKeys[key]; valid {
				next.ServeHTTP(w, r)
				return
			}
		
			http.Error(w, "API key not valid", http.StatusUnauthorized)
		})
	}
}

package auth

import (
	"net/http"
	"testing"

	"github.com/peterahlstrom/go-getter/config"
)

var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{"Success"})
})

// type Endpoint struct {
// 	UrlPath 		string `json:"urlPath"`
// 	ScriptPath 		string `json:"scriptPath"`
// 	RequireAuth		bool `json:"requireAuth"`
// 	ValidApiKeys	map[string]string `json:"apiKeys"`
// }

var tests = []map[string]config.Endpoint{


	{
		"no-auth": {
			ScriptPath: "",
			RequireAuth: false,
			ValidApiKeys: nil,
		},
		"validkey": {
			ScriptPath: "",
			RequireAuth: true,
			ValidApiKeys: map[string]string{"abc123": "test"},
		},
	},
}

func TestValidateApiKey(t *testing.T) {

	tests := []struct {
		name 	string
		apiKey	string
		endpoint map[string]config.Endpoint
		expectedStatus	int
	}{
		{"Valid API Key", "abc123", map[string]config.Endpoint{"endpoint": {config.Endpoint{ScriptPath: "", RequireAuth: true, ValidApiKeys: map[string]string{"abc123": "test"}}}}, expectedStatus: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ApiKeyMiddleWare(tt.endpoint, mockHandler)
		})
	}


}
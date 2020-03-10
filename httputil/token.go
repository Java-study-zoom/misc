package httputil

import (
	"net/http"
)

// SetAuthToken sets authorization header token.
func SetAuthToken(h http.Header, tok string) {
	if tok == "" {
		return
	}
	h.Set("Authorization", "Bearer "+tok)
}

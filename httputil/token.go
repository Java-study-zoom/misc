package httputil

import (
	"context"
	"net/http"
)

// SetAuthToken sets authorization header token.
func SetAuthToken(h http.Header, tok string) {
	if tok == "" {
		return
	}
	h.Set("Authorization", "Bearer "+tok)
}

// TokenSource is an interface that can provides a bearer token for
// authentication.
type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

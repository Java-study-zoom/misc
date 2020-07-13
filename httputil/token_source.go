package httputil

import (
	"context"
)

// TokenSource is an interface that can provides a bearer token for
// authentication.
type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

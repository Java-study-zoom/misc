package httputil

import (
	"net/http"
)

// AddToken adds the authorization header into the request.
func AddToken(req *http.Request, tok string) {
	headerSetAuthToken(req.Header, tok)
}

func headerSetAuthToken(h http.Header, tok string) {
	h.Set("Authorization", "Bearer "+tok)
}

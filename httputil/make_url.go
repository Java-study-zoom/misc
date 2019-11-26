package httputil

import (
	"net/url"
	"path"
)

func makeURL(server, p string) (string, error) {
	u, err := url.Parse(server)
	if err != nil {
		return "", err
	}

	up, err := url.Parse(p)
	if err != nil {
		return "", err
	}

	// append two paths
	u.Path = path.Join(u.Path, up.Path)
	u.RawQuery = up.RawQuery
	u.Fragment = up.Fragment
	return u.String(), nil
}

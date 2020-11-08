package jwt

import (
	"encoding/json"
)

// ClaimSet contains the JWT claims
type ClaimSet struct {
	Iss   string `json:"iss"`   // Issuer.
	Scope string `json:"scope"` // Scope, space-delimited list.
	Aud   string `json:"aud"`   // Audiance. Intended target.
	Exp   int64  `json:"exp"`   // Expiration time (Unix timestamp seconds)
	Iat   int64  `json:"iat"`   // Asserstion time (Unix timestamp seconds)
	Typ   string `json:"typ"`   // Token type.

	Sub string `json:"sub"`

	Extra map[string]interface{} `json:"-"`
}

func (c *ClaimSet) encode() (string, error) {
	m := make(map[string]interface{})

	for _, entry := range []struct {
		k, v     string
		mustHave bool
	}{
		{k: "iss", v: c.Iss, mustHave: true},
		{k: "scope", v: c.Scope},
		{k: "aud", v: c.Aud, mustHave: true},
		{k: "typ", v: c.Typ},
		{k: "sub", v: c.Sub},
	} {
		if entry.mustHave || entry.v != "" {
			m[entry.k] = entry.v
		}
	}

	m["iss"] = c.Iss
	m["iat"] = c.Iat

	for k, v := range c.Extra {
		m[k] = v
	}

	return encodeSegment(m)
}

func decodeClaimSet(bs []byte) (*ClaimSet, error) {
	c := new(ClaimSet)
	if err := json.Unmarshal(bs, c); err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}

	for _, k := range []string{
		"iss", "scope", "aud", "typ", "sub", "iss", "iat",
	} {
		delete(m, k)
	}
	if len(m) > 0 {
		c.Extra = m
	}
	return c, nil
}

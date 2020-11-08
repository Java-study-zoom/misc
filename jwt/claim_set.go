package jwt

// ClaimSet contains the JWT claims
type ClaimSet struct {
	Iss   string // Issuer.
	Scope string // Scope, space-delimited list.
	Aud   string // Audiance. Intended target.
	Exp   int64  // Expiration time (Unix timestamp seconds)
	Iat   int64  // Asserstion time (Unix timestamp seconds)
	Typ   string // Token type.

	Sub string

	Extra map[string]interface{}
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

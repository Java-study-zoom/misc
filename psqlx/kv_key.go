package psqlx

import (
	"shanhu.io/misc/hashutil"
)

func keyHash(k string) string {
	return hashutil.HashStr(k)
}

// MaxKeyLen is the maximum length of a hashed KV.
const MaxKeyLen = 255

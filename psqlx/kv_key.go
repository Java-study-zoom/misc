package psqlx

import (
	"database/sql"
	"errors"
	"fmt"

	"shanhu.io/misc/hashutil"
	"shanhu.io/misc/pathutil"
)

// MaxKeyLen is the maximum length of a hashed KV.
const MaxKeyLen = 255

func keyHash(k string) string {
	return hashutil.HashStr(k)
}

func kvMapKey(key string, hashed bool) (string, error) {
	if hashed {
		return keyHash(key), nil
	}
	if len(key) > MaxKeyLen {
		return "", fmt.Errorf("key %q too long", key)
	}
	return key, nil
}

func kvResError(res sql.Result, key string) error {
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return pathutil.NotExist(key)
	} else if n == 1 {
		return errors.New("multiple rows affected")
	}
	return nil
}

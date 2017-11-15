package psqlx

import (
	"database/sql"
	"encoding/json"
)

func kvIterRows(rows *sql.Rows, it *Iter, hashed bool) error {
	for rows.Next() {
		var k string
		var bs []byte
		if err := rows.Scan(&k, &bs); err != nil {
			return err
		}

		entry := it.Make()
		if err := json.Unmarshal(bs, entry); err != nil {
			return err
		}
		if hashed {
			k = ""
		}
		if err := it.Do(k, entry); err != nil {
			return err
		}
	}

	return rows.Close()
}

func orderStr(desc bool) string {
	if desc {
		return "desc"
	}
	return "acs"
}

// KVPartial specifies a partial of results.
type KVPartial struct {
	Offset uint64
	N      uint64
	Desc   bool
}

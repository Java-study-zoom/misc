package psqlx

import (
	"fmt"

	"shanhu.io/misc/sqlx"
)

// InitKV creates a key value pair
func InitKV(db *sqlx.DB, table string) error {
	q := fmt.Sprintf(`create table %s (
		k varchar(%d) primary key not null,
		v bytea not null
	)`, table, MaxKeyLen)
	_, err := db.X(q)
	return err
}

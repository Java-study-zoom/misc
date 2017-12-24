package psqlx

import (
	"fmt"

	"shanhu.io/misc/sqlx"
)

// MaxClassLen is the maximum length of the class string of a hashed KV.
const MaxClassLen = 255

// InitKV creates a key value pair table.
func InitKV(db *sqlx.DB, table string) error {
	q := fmt.Sprintf(`create table %s (
		k varchar(%d) primary key not null,
		c varchar(%d) not null,
		v bytea not null
	)`, table, MaxKeyLen, MaxClassLen)
	_, err := db.X(q)
	return err
}

// InitKVs creates a series of key value pair tables.
func InitKVs(db *sqlx.DB, tables ...string) error {
	for _, table := range tables {
		if err := InitKV(db, table); err != nil {
			return err
		}
	}
	return nil
}

// KVAddClassColumn adds a class column for the KV table.
func KVAddClassColumn(db *sqlx.DB, table string) error {
	q := fmt.Sprintf(
		`alter table %s add column c varchar(%d) not null default ''`,
		table, MaxClassLen,
	)
	_, err := db.X(q)
	return err
}

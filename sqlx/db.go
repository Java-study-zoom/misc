package sqlx

import (
	"database/sql"
)

// DB is a wrapper that extends the sql.DB structure.
type DB struct {
	*sql.DB
	*wrap
}

// OpenPsql opens a postgresql database.
func OpenPsql(source string) (*DB, error) {
	if source == "" {
		return nil, nil
	}
	return Open("postgres", source)
}

// Open opens a database.
func Open(driver, source string) (*DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:   db,
		wrap: &wrap{conn: db},
	}, nil
}

// Begin begins a transaction.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:   tx,
		wrap: &wrap{conn: tx},
	}, nil
}

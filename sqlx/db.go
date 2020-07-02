package sqlx

import (
	"database/sql"
)

// DB is a wrapper that extends the sql.DB structure.
type DB struct {
	*sql.DB
	*wrap

	driver string
}

// Driver names.
const (
	Psql    = "postgres"
	Sqlite3 = "sqlite3"
)

// OpenPsql opens a postgresql database.
func OpenPsql(source string) (*DB, error) {
	if source == "" {
		return nil, nil
	}
	return Open(Psql, source)
}

// OpenSqlite3 opens a sqlite3 database.
func OpenSqlite3(file string) (*DB, error) {
	if file == "" {
		return nil, nil
	}
	return Open(Sqlite3, file)
}

// Open opens a database.
func Open(driver, source string) (*DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:     db,
		wrap:   &wrap{conn: db},
		driver: driver,
	}, nil
}

// Driver returns the driver name when the database is being opened.
func (db *DB) Driver() string {
	return db.driver
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

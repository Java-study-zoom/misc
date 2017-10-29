package psqlx

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"shanhu.io/misc/pathutil"
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

// KV defines a key value pair table.
type KV struct {
	db     *sqlx.DB
	table  string
	hashed bool
}

// NewKV connects to a generic unordered key-value pair table.
func NewKV(db *sqlx.DB, table string) *KV {
	return &KV{db: db, table: table, hashed: true}
}

// NewOrderedKV connects to a generic ordered key-value pair table.
func NewOrderedKV(db *sqlx.DB, table string) *KV {
	return &KV{db: db, table: table, hashed: false}
}

func (b *KV) mapKey(key string) (string, error) {
	if b.hashed {
		return keyHash(key), nil
	}
	if len(key) > MaxKeyLen {
		return "", fmt.Errorf("key %q too long", key)
	}
	return key, nil
}

// Add adds an entry with the given key and value. The value
// will be marshalled with JSON encoding.
func (b *KV) Add(key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`insert into %s (k, v) values ($1, $2)`, b.table)
	_, err = b.db.X(q, mk, bs)
	return err
}

// Remove removes the entry with the specific key.
func (b *KV) Remove(key string) error {
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`delete from %s where k=$1`, b.table)
	res, err := b.db.X(q, mk)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return pathutil.NotExist(key)
	}
	return nil
}

// GetBytes gets the value bytes for the specific key.
func (b *KV) GetBytes(key string) ([]byte, error) {
	mk, err := b.mapKey(key)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`select v from %s where k=$1`, b.table)
	row := b.db.Q1(q, mk)
	var bs []byte
	if has, err := row.Scan(&bs); err != nil {
		return nil, err
	} else if !has {
		return nil, pathutil.NotExist(key)
	}
	return bs, nil
}

// Get gets the value and json marshals it into v.
func (b *KV) Get(key string, v interface{}) error {
	bs, err := b.GetBytes(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// Emplace sets the value for a particular key. Creates the key if not exist.
func (b *KV) Emplace(key string, v interface{}) error {
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		insert into %s (k, v) values ($1, $2)
		on conflict (k) do update set v=excluded.v
	`, b.table)
	_, err = b.db.X(q, mk, bs)
	return err
}

// AppendBytes appends the value to the existing value of a particular key.
// Creates the key if not exist.
func (b *KV) AppendBytes(key string, bs []byte) error {
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		insert into %s (k, v) values ($1, $2)
		on conflict (k) do update set v = v || excluded.v
	`, b.table)
	_, err = b.db.X(q, mk, bs)
	return err
}

// SetBytes updates the value bytes of a particular key.
func (b *KV) SetBytes(key string, bs []byte) error {
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`update %s set v=$1 where k=$2`, b.table)
	res, err := b.db.X(q, bs, mk)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return pathutil.NotExist(key)
	}
	if n != 1 {
		return errors.New("multiple value got updated")
	}
	return nil
}

// Set updates the JSON value of a particular key.
func (b *KV) Set(key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.SetBytes(key, bs)
}

// ErrCancel cancels the operation.
var ErrCancel = errors.New("operation cancelled")

// Mutate applies a function
func (b *KV) Mutate(
	k string, v interface{}, f func(v interface{}) error,
) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}

	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bs []byte
	q := fmt.Sprintf(`select v from %s where k=$1`, b.table)
	row := tx.Q1(q, mk)
	if has, err := row.Scan(&bs); err != nil {
		return err
	} else if !has {
		return pathutil.NotExist(k)
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}

	err = f(v)
	if err == ErrCancel {
		return nil
	}
	if err != nil {
		return err
	}

	bs, err = json.Marshal(v)
	if err != nil {
		return err
	}

	q = fmt.Sprintf(`update %s set v=$1 where k=$2`, b.table)
	res, err := tx.X(q, bs, mk)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("nothing updated")
	}
	if n != 1 {
		return fmt.Errorf("%d updated", n)
	}
	return tx.Commit()
}

func iterRows(rows *sql.Rows, it *Iter) error {
	for rows.Next() {
		var bs []byte
		if err := rows.Scan(&bs); err != nil {
			return err
		}

		entry := it.Make()
		if err := json.Unmarshal(bs, entry); err != nil {
			return err
		}
		if err := it.Do(entry); err != nil {
			return err
		}
	}

	return rows.Close()
}

// Walk iterates through all items in the key value store.
func (b *KV) Walk(it *Iter) error {
	q := fmt.Sprintf(`select k, v from %s order by k`, b.table)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	return iterRows(rows, it)
}

// WalkPartial walks thorugh the items at offset with at most n items.
func (b *KV) WalkPartial(offset, n uint64, desc bool, it *Iter) error {
	if b.hashed {
		return fmt.Errorf("cannot partial walk over a hashed table")
	}
	q := fmt.Sprintf(
		`select k, v from %s order by k offset ? limit ?`, b.table,
	)
	if desc {
		q += " desc"
	}
	rows, err := b.db.Q(q, offset, n)
	if err != nil {
		return err
	}
	defer rows.Close()

	return iterRows(rows, it)
}

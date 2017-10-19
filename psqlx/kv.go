package psqlx

import (
	"encoding/json"
	"errors"
	"fmt"

	"shanhu.io/misc/hashutil"
	"shanhu.io/misc/pathutil"
	"shanhu.io/misc/sqlx"
)

func keyHash(k string) string {
	return hashutil.HashStr(k)
}

// InitKV creates a key value pair
func InitKV(db *sqlx.DB, table string) error {
	q := fmt.Sprintf(`create table %s (
		k varchar(255) primary key not null,
		v bytea not null
	)`, table)
	_, err := db.X(q)
	return err
}

// KV defines a key value table.
type KV struct {
	db    *sqlx.DB
	table string
}

// NewKV is a generic key-value pair table.
func NewKV(db *sqlx.DB, table string) *KV {
	return &KV{db: db, table: table}
}

// Add adds an entry with the given key and value. The value
// will be marshalled with JSON encoding.
func (b *KV) Add(key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`insert into %s (k, v) values ($1, $2)`, b.table)
	_, err = b.db.X(q, keyHash(key), bs)
	return err
}

// Remove removes the entry with the specific key.
func (b *KV) Remove(key string) error {
	q := fmt.Sprintf(`delete from %s where k=$1`, b.table)
	res, err := b.db.X(q, keyHash(key))
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

// Get gets the value and json marshals it into v.
func (b *KV) Get(key string, v interface{}) error {
	q := fmt.Sprintf(`select v from %s where k=$1`, b.table)
	row := b.db.Q1(q, keyHash(key))
	var bs []byte
	if has, err := row.Scan(&bs); err != nil {
		return err
	} else if !has {
		return pathutil.NotExist(key)
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}
	return nil
}

// Emplace sets the value for a particular key. Creates the key if not exist.
func (b *KV) Emplace(key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		insert into %s (k, v) values ($1, $2)
		on conflict (k) do update set v=excluded.v
	`, b.table)
	_, err = b.db.X(q, keyHash(key), bs)
	return err
}

// Set updates a value of the particular key.
func (b *KV) Set(key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`update %s set v=$1 where k=$2`, b.table)
	res, err := b.db.X(q, bs, keyHash(key))
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

// ErrCancel cancels the operation.
var ErrCancel = errors.New("operation cancelled")

// Mutate applies a function
func (b *KV) Mutate(
	k string, v interface{}, f func(v interface{}) error,
) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	hash := keyHash(k)

	var bs []byte
	q := fmt.Sprintf(`select v from %s where k=$1`, b.table)
	row := tx.Q1(q, hash)
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
	res, err := tx.X(q, bs, hash)
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

// Walk iterates through all items in the key value store.
func (b *KV) Walk(
	makeFunc func() interface{}, f func(v interface{}) error,
) error {
	q := fmt.Sprintf(`select k, v from %s order by k`, b.table)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var bs []byte
		if err := rows.Scan(&bs); err != nil {
			return err
		}

		entry := makeFunc()
		if err := json.Unmarshal(bs, entry); err != nil {
			return err
		}
		if err := f(entry); err != nil {
			return err
		}
	}

	return rows.Close()
}

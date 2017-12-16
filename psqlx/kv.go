package psqlx

import (
	"encoding/json"
	"errors"
	"fmt"

	"shanhu.io/misc/pathutil"
	"shanhu.io/misc/sqlx"
)

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

func (b *KV) mapKey(k string) (string, error) {
	return kvMapKey(k, b.hashed)
}

// AddClass adds an entry with a particular class.
func (b *KV) AddClass(key, cls string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`insert into %s (k, c, v) values ($1, $2, $3)`, b.table)
	_, err = b.db.X(q, mk, cls, bs)
	return err
}

// SetClass sets an entry's class.
func (b *KV) SetClass(key, cls string) error {
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`update %s set c=$1 where k=$2`, b.table)
	res, err := b.db.X(q, cls, mk)
	if err != nil {
		return err
	}
	return kvResError(res, key)
}

// Add adds an entry with the given key and value. The value
// will be marshalled with JSON encoding.
func (b *KV) Add(key string, v interface{}) error {
	return b.AddClass(key, "", v)
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

	return kvResError(res, key)
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

// Emplace sets the value for a particular key. Does nothing if the key does
// not exist.
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
        insert into %s (k, v, c) values ($1, $2, $3)
        on conflict (k) do nothing
    `, b.table)
	_, err = b.db.X(q, mk, bs, "")
	return err
}

// Replace sets the value for a particular key. Creates the key if not exist.
func (b *KV) Replace(key string, v interface{}) error {
	mk, err := b.mapKey(key)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		insert into %s (k, v, c) values ($1, $2, $3)
		on conflict (k) do update set v=excluded.v
	`, b.table)
	_, err = b.db.X(q, mk, bs, "")
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
		insert into %s (k, v, c) values ($1, $2, $3)
		on conflict (k) do update set v = %s.v || excluded.v
	`, b.table, b.table)
	_, err = b.db.X(q, mk, bs, "")
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

	return kvResError(res, key)
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

// Mutate applies a function to an item.
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

	if err := f(v); err == ErrCancel {
		return nil
	} else if err != nil {
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
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return fmt.Errorf("nothing updated")
	} else if n != 1 {
		return fmt.Errorf("%d updated", n)
	}
	return tx.Commit()
}

// Walk iterates through all items in the key value store.
func (b *KV) Walk(it *Iter) error {
	q := fmt.Sprintf(`select k, v from %s order by k`, b.table)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return kvIterRows(rows, it, b.hashed)
}

var errHashedHasNoPartial = fmt.Errorf(
	"cannot partial walk over a hashed table",
)

// WalkPartial walks through the some part of the resulting items.
func (b *KV) WalkPartial(p *KVPartial, it *Iter) error {
	if b.hashed {
		return errHashedHasNoPartial
	}
	q := fmt.Sprintf(
		`select k, v from %s order by k %s limit %d offset %d`,
		b.table, orderStr(p.Desc), p.N, p.Offset,
	)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return kvIterRows(rows, it, b.hashed)
}

// WalkPartialClass walks through the some part of the items that
// of a given class.
func (b *KV) WalkPartialClass(cls string, p *KVPartial, it *Iter) error {
	if b.hashed {
		return errHashedHasNoPartial
	}
	q := fmt.Sprintf(
		`select k, v from %s where c=$1
		order by k %s limit %d offset %d`,
		b.table, orderStr(p.Desc), p.N, p.Offset,
	)
	rows, err := b.db.Q(q, cls)
	if err != nil {
		return err
	}
	defer rows.Close()
	return kvIterRows(rows, it, b.hashed)
}

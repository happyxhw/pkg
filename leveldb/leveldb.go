package leveldb

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB(path string, o *opt.Options) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}
	t := &LevelDB{
		db: db,
	}

	return t, nil
}

func (t *LevelDB) Get(key []byte, ro *opt.ReadOptions) ([]byte, error) {
	return t.db.Get(key, ro)
}

func (t *LevelDB) GetObject(key []byte, out any, ro *opt.ReadOptions) error {
	data, err := t.db.Get(key, ro)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &out)
}

func (t *LevelDB) Put(key, value []byte, wo *opt.WriteOptions) error {
	return t.db.Put(key, value, wo)
}

func (t *LevelDB) PutObject(key []byte, value any, wo *opt.WriteOptions) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return t.db.Put(key, data, wo)
}

func (t *LevelDB) Del(key []byte, wo *opt.WriteOptions) error {
	return t.db.Delete(key, wo)
}

func (t *LevelDB) Iter(fn func(key, value []byte), slice *util.Range, ro *opt.ReadOptions) error {
	iter := t.db.NewIterator(slice, ro)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fn(key, value)
	}
	iter.Release()
	return iter.Error()
}

func (t *LevelDB) Close() error {
	return t.db.Close()
}

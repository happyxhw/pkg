package queue

import (
	"encoding/binary"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	readPositionFlag  = []byte("_readPosition")
	writePositionFlag = []byte("_writePosition")
)

type LevelQueue struct {
	db            *leveldb.DB
	mutex         sync.Mutex
	cond          *sync.Cond
	readPosition  uint64
	writePosition uint64
	capacity      int
}

func NewLevelQueue(dbPath string, capacity int) (*LevelQueue, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	var readPosition uint64 = 0
	var writePosition uint64 = 0

	readBufs, err := db.Get(readPositionFlag, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	} else if readBufs != nil {
		readPosition = binary.BigEndian.Uint64(readBufs)
	}
	writeBufs, err := db.Get(writePositionFlag, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	} else if writeBufs != nil {
		writePosition = binary.BigEndian.Uint64(writeBufs)
	}

	q := LevelQueue{
		db:            db,
		readPosition:  readPosition,
		writePosition: writePosition,
		capacity:      capacity,
	}
	q.cond = sync.NewCond(&q.mutex)
	return &q, nil
}

func (q *LevelQueue) Push(data []byte) (bool, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.capacity > 0 {
		for q.writePosition-q.readPosition >= uint64(q.capacity) {
			q.cond.Wait()
		}
	}

	pos := make([]byte, 8)
	binary.BigEndian.PutUint64(pos, q.writePosition)
	if err := q.db.Put(pos, data, nil); err != nil {
		return false, err
	}

	binary.BigEndian.PutUint64(pos, q.writePosition+1)
	if err := q.db.Put(writePositionFlag, pos, nil); err != nil {
		return false, err
	}
	q.writePosition += 1
	q.cond.Signal()
	return true, nil
}

// Pop: Last + DeleteLast
func (q *LevelQueue) Pop() ([]byte, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.readPosition >= q.writePosition {
		q.readPosition = 0
		q.writePosition = 0
		q.cond.Wait()
	}
	pos := make([]byte, 8)
	binary.BigEndian.PutUint64(pos, q.readPosition)
	value, err := q.db.Get(pos, nil)
	if err != nil && err != leveldb.ErrNotFound { // error
		return nil, err
	}
	if value != nil { // exists
		if dErr := q.db.Delete(pos, nil); dErr != nil {
			return nil, dErr
		}
	}
	binary.BigEndian.PutUint64(pos, q.readPosition+1)
	err = q.db.Put(readPositionFlag, pos, nil)
	if err != nil {
		return nil, err
	}
	q.readPosition += 1
	q.cond.Signal()
	return value, nil
}

func (q *LevelQueue) Last() ([]byte, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.readPosition >= q.writePosition {
		q.readPosition = 0
		q.writePosition = 0
		q.cond.Wait()
	}
	pos := make([]byte, 8)
	binary.BigEndian.PutUint64(pos, q.readPosition)
	value, err := q.db.Get(pos, nil)
	if err != nil && err != leveldb.ErrNotFound { // error
		return nil, err
	}
	return value, nil
}

func (q *LevelQueue) DeleteLast() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	pos := make([]byte, 8)
	binary.BigEndian.PutUint64(pos, q.readPosition)
	if err := q.db.Delete(pos, nil); err != nil {
		return err
	}
	binary.BigEndian.PutUint64(pos, q.readPosition+1)
	err := q.db.Put(readPositionFlag, pos, nil)
	if err != nil {
		return err
	}
	q.readPosition += 1
	q.cond.Signal()
	return nil
}

func (q *LevelQueue) DestroyQueue() {
	q.db.Close()
}

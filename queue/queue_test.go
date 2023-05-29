package queue

import (
	"strconv"
	"testing"
	"time"
)

func TestNewLevelQueue(t *testing.T) {
	q, err := NewLevelQueue("/tmp/database")
	if err != nil {
		panic(err)
	}

	go func() {
		var i int
		for {
			_, err := q.Push([]byte(strconv.Itoa(i)))
			if err != nil {
				panic(err)
			}
			i++
			time.Sleep(time.Millisecond * 100)
			if i > 10 {
				break
			}
		}
	}()

	for {
		item, err := q.Pop()
		if err != nil {
			t.Fatal(err)
		}

		data, _ := strconv.Atoi(string(item))
		if data >= 10 {
			break
		}
	}
}

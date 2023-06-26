package queue

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNewLevelQueue(t *testing.T) {
	q, err := NewLevelQueue("/tmp/database", 10)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll("/tmp/database")

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
		fmt.Println(data)
		if data >= 10 {
			break
		}
	}
}

func TestNewLevelQueueLast(t *testing.T) {
	q, err := NewLevelQueue("/tmp/database", 10)
	if err != nil {
		panic(err)
	}
	// defer os.RemoveAll("/tmp/database")
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
		item, err := q.Last()
		if err != nil {
			t.Fatal(err)
		}
		if err := q.DeleteLast(); err != nil {
			t.Fatal(err)
		}

		data, _ := strconv.Atoi(string(item))
		fmt.Println(data)
		if data >= 10 {
			break
		}
	}
}

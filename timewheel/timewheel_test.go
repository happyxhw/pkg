package timewheel

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New(time.Millisecond*100, 100, func(data interface{}) {})
	if err != nil {
		t.Error(err)
	}
	type args struct {
		interval time.Duration
		slotNum  int
		job      Job
	}
	tests := []struct {
		name    string
		args    args
		want    *TimeWheel
		wantErr bool
	}{
		{
			name:    "invalid interval",
			args:    args{interval: time.Duration(-1), slotNum: 100, job: func(data interface{}) {}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid slot num",
			args:    args{interval: time.Millisecond * 100, slotNum: -1, job: func(data interface{}) {}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.interval, tt.args.slotNum, tt.args.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTimer(t *testing.T) {
	key := "task-1"
	jobData := map[string]string{"k": "v"}

	start := time.Now()
	delay := time.Millisecond * 500

	tw, err := New(time.Millisecond*10, 500, func(data interface{}) {
		elapse := time.Since(start).Milliseconds()

		require.LessOrEqual(t, elapse-delay.Milliseconds(), int64(20))
		require.Equal(t, data, jobData)
	})
	if err != nil {
		t.Error(err)
	}
	tw.Start()

	tw.AddTimer(delay, key, jobData)

	time.Sleep(time.Second)
}

func TestRemoveTimer(t *testing.T) {
	key := "task-1"
	jobData := map[string]string{"k": "v"}
	delay := time.Millisecond * 500

	tw, err := New(time.Millisecond*10, 500, func(data interface{}) {
		t.Error("should not be invoked")
	})
	if err != nil {
		t.Error(err)
	}
	tw.Start()

	tw.AddTimer(delay, key, jobData)
	tw.RemoveTimer(key)
	if _, ok := tw.tasks[key]; ok {
		t.Error("timer still exists")
	}

	time.Sleep(time.Second)
}

func TestUpdateTimer(t *testing.T) {
	key := "task-1"
	jobData := map[string]string{"k": "v"}
	delay1 := time.Millisecond * 100
	delay2 := time.Millisecond * 500
	start := time.Now()
	tw, err := New(time.Millisecond*10, 500, func(data interface{}) {
		elapse := time.Since(start).Milliseconds()
		require.LessOrEqual(t, elapse-delay2.Milliseconds(), int64(20))
		require.Equal(t, data, jobData)
	})
	if err != nil {
		t.Error(err)
	}
	tw.Start()

	tw.AddTimer(delay1, key, jobData)
	start = time.Now()
	tw.UpdateTimer(delay2, key, jobData)

	time.Sleep(time.Second)
}

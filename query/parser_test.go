package query

import (
	"testing"
)

func TestParseOrder(t *testing.T) {
	type args struct {
		sortBy  string
		handles []FieldFunc
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{sortBy: "+created_at"},
			want: "created_at ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseOrder(tt.args.sortBy, tt.args.handles...); got != tt.want {
				t.Errorf("ParseOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

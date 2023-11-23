package query

import (
	"reflect"
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

func TestParseFilter(t *testing.T) {
	type args struct {
		filter  string
		handles []FilterFunc
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []any
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				// or(and(eq(name,h),le(age,10),or(eq(name,y),le(age,11))),or(in(status,0,1,2),nin(role,1,2,3)))
				filter: `
					or
					(
						and(
							eq(name,h),
							le(age,10),
							or
							(
								eq(name,y),
								le(age,11)
							)
						),
						or
						(
							in(status,0,1,2),
							nin(role,1,2,3)
						)
					)
				`,
			},
			want:    "((role not in ? or status in ?) or ((age <= ? or name = ?) and age <= ? and name = ?))",
			want1:   []any{[]any{"1", "2", "3"}, []any{"0", "1", "2"}, "11", "y", "10", "h"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseFilter(tt.args.filter, tt.args.handles...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

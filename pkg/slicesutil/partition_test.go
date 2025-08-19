package slicesutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Partition(t *testing.T) {
	type args struct {
		values []int
		size   int
	}
	type want struct {
		result [][]int
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "with empty values",
			args: &args{
				values: []int{},
				size:   10,
			},
			want: &want{
				result: [][]int{},
			},
		},
		{
			name: "with small values length",
			args: &args{
				values: []int{1, 2},
				size:   10,
			},
			want: &want{
				result: [][]int{
					{1, 2},
				},
			},
		},
		{
			name: "with exact values length as size",
			args: &args{
				values: []int{1, 2},
				size:   2,
			},
			want: &want{
				result: [][]int{
					{1, 2},
				},
			},
		},
		{
			name: "with greater values length as size",
			args: &args{
				values: []int{1, 2, 3, 4, 5},
				size:   2,
			},
			want: &want{
				result: [][]int{
					{1, 2},
					{3, 4},
					{5},
				},
			},
		},
		{
			name: "with greater values length as size",
			args: &args{
				values: []int{1, 2, 3, 4, 5, 6},
				size:   3,
			},
			want: &want{
				result: [][]int{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Partition(tt.args.values, tt.args.size)

			assert.Equal(t, tt.want.result, result)
		})
	}
}

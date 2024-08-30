package asyncutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
	"github.com/vitorsss/go-helpers/pkg/asyncutil"
)

func Test_ConcurrencyExec(t *testing.T) {
	type args struct {
		items   []int
		fn      func(ctx context.Context, item int) (int, error)
		options []asyncutil.ConcurrencyExecOption
	}
	type want struct {
		result       []int
		expectedTime time.Duration
		err          error
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should execute concurrently",
			args: &args{
				items: []int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				fn: func(ctx context.Context, item int) (int, error) {
					time.Sleep(10 * time.Millisecond)
					return item, nil
				},
			},
			want: &want{
				result: []int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				expectedTime: 30 * time.Millisecond,
			},
		},
		{
			name: "should execute concurrently and return the nested error",
			args: &args{
				items: []int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				fn: func(ctx context.Context, item int) (int, error) {
					time.Sleep(10 * time.Millisecond)
					if item == 20 {
						return 0, errors.New("some nested error")
					}
					return item, nil
				},
			},
			want: &want{
				result:       nil,
				expectedTime: 30 * time.Millisecond,
				err:          errors.New("some nested error"),
			},
		},
		{
			name: "should execute concurrently and return the nested panic with partial success",
			args: &args{
				items: []int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				fn: func(ctx context.Context, item int) (int, error) {
					time.Sleep(10 * time.Millisecond)
					if item == 2 {
						panic("some nested panic")
					}
					return item, nil
				},
				options: []asyncutil.ConcurrencyExecOption{
					asyncutil.WithWaitPartialSuccess(),
				},
			},
			want: &want{
				result: []int{
					1, 3, 4, 5, 6, 7,
				},
				expectedTime: 30 * time.Millisecond,
				err:          errors.New("recovered error on safe exec: some nested panic"),
			},
		},
		{
			name: "should execute concurrently with custom max concurrency",
			args: &args{
				items: []int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				fn: func(ctx context.Context, item int) (int, error) {
					time.Sleep(10 * time.Millisecond)
					return item, nil
				},
				options: []asyncutil.ConcurrencyExecOption{
					asyncutil.WithMaxConcurrency(5),
				},
			},
			want: &want{
				result: []int{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				expectedTime: 50 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			result, err := asyncutil.ConcurrencyExec(
				context.Background(),
				tt.args.items,
				tt.args.fn,
				tt.args.options...,
			)
			time := time.Now().Sub(start)

			if assertutil.Error(t, tt.want.err, err) {
				for _, expected := range tt.want.result {
					assert.Contains(t, result, expected)
				}
				assert.LessOrEqual(t, time, tt.want.expectedTime)
			}
		})
	}
}

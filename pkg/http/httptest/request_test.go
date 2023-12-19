package httptest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_matchURLValues(t *testing.T) {
	type args struct {
		expected url.Values
		current  url.Values
	}

	type want struct {
		result bool
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should match successfully",
			args: &args{
				expected: url.Values{
					"param": []string{"value"},
				},
				current: url.Values{
					"param": []string{"value"},
				},
			},
			want: &want{
				result: true,
			},
		},
		{
			name: "should match successfully unordered params",
			args: &args{
				expected: url.Values{
					"param": []string{"value", "valueb"},
				},
				current: url.Values{
					"param": []string{"valueb", "value"},
				},
			},
			want: &want{
				result: true,
			},
		},
		{
			name: "should not match since missing total params",
			args: &args{
				expected: url.Values{
					"param": []string{"value"},
				},
				current: url.Values{},
			},
			want: &want{
				result: false,
			},
		},
		{
			name: "should not match since too many total params",
			args: &args{
				expected: url.Values{},
				current: url.Values{
					"param": []string{"value"},
				},
			},
			want: &want{
				result: false,
			},
		},
		{
			name: "should not match since missing param",
			args: &args{
				expected: url.Values{
					"param": []string{"value"},
				},
				current: url.Values{
					"paramb": []string{"value"},
				},
			},
			want: &want{
				result: false,
			},
		},
		{
			name: "should not match since param values does not match",
			args: &args{
				expected: url.Values{
					"param": []string{"value"},
				},
				current: url.Values{
					"param": []string{"valuea"},
				},
			},
			want: &want{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchURLValues(tt.args.expected, tt.args.current)

			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_matchHeader(t *testing.T) {
	type args struct {
		expected http.Header
		current  http.Header
	}

	type want struct {
		result bool
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should match successfully",
			args: &args{
				expected: http.Header{
					"Param": []string{"value"},
				},
				current: http.Header{
					"Param": []string{"value"},
				},
			},
			want: &want{
				result: true,
			},
		},
		{
			name: "should match successfully unordered headers",
			args: &args{
				expected: http.Header{
					"Param": []string{"value", "valueb"},
				},
				current: http.Header{
					"Param": []string{"valueb", "value"},
				},
			},
			want: &want{
				result: true,
			},
		},
		{
			name: "should match successfully with extra header",
			args: &args{
				expected: http.Header{
					"Param": []string{"value", "valueb"},
				},
				current: http.Header{
					"Param":  []string{"valueb", "value"},
					"Paramb": []string{"valueb", "value"},
				},
			},
			want: &want{
				result: true,
			},
		},
		{
			name: "should not match since header values does not match",
			args: &args{
				expected: http.Header{
					"Param": []string{"value"},
				},
				current: http.Header{
					"Param": []string{"valuea"},
				},
			},
			want: &want{
				result: false,
			},
		},
		{
			name: "should not match since missing header",
			args: &args{
				expected: http.Header{
					"Param": []string{"value"},
				},
				current: http.Header{},
			},
			want: &want{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchHeader(tt.args.expected, tt.args.current)

			assert.Equal(t, tt.want.result, result)
		})
	}
}

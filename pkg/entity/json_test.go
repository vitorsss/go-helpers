package entity_test

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
	"github.com/vitorsss/go-helpers/pkg/entity"
)

type testEntity struct {
	ID int `json:"id"`
}

func Test_JSONStringWrapper_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}

	type want struct {
		value entity.JSONStringWrapper[testEntity]
		err   error
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should unmarshal content successfully",
			args: &args{
				data: []byte(`"{\"id\":10}"`),
			},
			want: &want{
				value: entity.JSONStringWrapper[testEntity]{
					Content: testEntity{
						ID: 10,
					},
				},
				err: nil,
			},
		},
		{
			name: "should return error while unmarshalling content",
			args: &args{
				data: []byte(`"{\"id\":\"10\"}"`),
			},
			want: &want{
				value: entity.JSONStringWrapper[testEntity]{},
				err:   errors.New("failed to unmarshal content: json: cannot unmarshal string into Go struct field testEntity.id of type int"),
			},
		},
		{
			name: "should return error while unmarshalling string",
			args: &args{
				data: []byte(`10`),
			},
			want: &want{
				value: entity.JSONStringWrapper[testEntity]{},
				err:   errors.New("failed to unmarshal data string: json: cannot unmarshal number into Go value of type string"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value entity.JSONStringWrapper[testEntity]
			err := json.Unmarshal(tt.args.data, &value)

			assertutil.Error(t, tt.want.err, err)
			assert.Equal(t, tt.want.value, value)
		})
	}
}

func Test_JSONStringWrapper_MarshalJSON(t *testing.T) {
	type args struct {
		value entity.JSONStringWrapper[any]
	}

	type want struct {
		data string
		err  error
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should marshal successfully",
			args: &args{
				value: entity.JSONStringWrapper[any]{
					Content: testEntity{
						ID: 10,
					},
				},
			},
			want: &want{
				data: `"{\"id\":10}"`,
				err:  nil,
			},
		},
		{
			name: "should return error while marshalling content",
			args: &args{
				value: entity.JSONStringWrapper[any]{
					Content: make(chan int),
				},
			},
			want: &want{
				data: ``,
				err:  errors.New("json: error calling MarshalJSON for type entity.JSONStringWrapper[interface {}]: failed to marshal content: json: unsupported type: chan int"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.args.value)

			assertutil.Error(t, tt.want.err, err)
			assert.Equal(t, tt.want.data, string(data))
		})
	}
}

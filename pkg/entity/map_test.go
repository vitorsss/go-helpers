package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/entity"
)

func Test_Map_Get(t *testing.T) {
	type args struct {
		mapValue entity.Map
		keys     []string
	}

	type want struct {
		value interface{}
		ok    bool
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should return nil value when empty keys",
			args: &args{
				mapValue: entity.Map{},
				keys:     []string{},
			},
			want: &want{
				value: nil,
				ok:    false,
			},
		},
		{
			name: "should return nil value when missing key value",
			args: &args{
				mapValue: entity.Map{},
				keys: []string{
					"key_a",
				},
			},
			want: &want{
				value: nil,
				ok:    false,
			},
		},
		{
			name: "should return nil value when missing keys values",
			args: &args{
				mapValue: entity.Map{},
				keys: []string{
					"key_a",
					"key_b",
				},
			},
			want: &want{
				value: nil,
				ok:    false,
			},
		},
		{
			name: "should return nil value when a intermediary key is not a map",
			args: &args{
				mapValue: entity.Map{
					"key_a": "test",
				},
				keys: []string{
					"key_a",
					"key_b",
				},
			},
			want: &want{
				value: nil,
				ok:    false,
			},
		},
		{
			name: "should return value when keys match structure of map[string]interface{}",
			args: &args{
				mapValue: entity.Map{
					"key_a": map[string]interface{}{
						"key_b": "test",
					},
				},
				keys: []string{
					"key_a",
					"key_b",
				},
			},
			want: &want{
				value: "test",
				ok:    true,
			},
		},
		{
			name: "should return value when keys match structure of Map",
			args: &args{
				mapValue: entity.Map{
					"key_a": entity.Map{
						"key_b": "test",
					},
				},
				keys: []string{
					"key_a",
					"key_b",
				},
			},
			want: &want{
				value: "test",
				ok:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := tt.args.mapValue.Get(tt.args.keys...)

			assert.Equal(t, tt.want.value, value)
			assert.Equal(t, tt.want.ok, ok)
		})
	}
}

func Test_Map_Set(t *testing.T) {
	type args struct {
		mapValue entity.Map
		value    interface{}
		keys     []string
	}

	type want struct {
		mapValue entity.Map
		updated  bool
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should not update when empty keys",
			args: &args{
				mapValue: entity.Map{},
				value:    "test",
				keys:     []string{},
			},
			want: &want{
				mapValue: entity.Map{},
				updated:  false,
			},
		},
		{
			name: "should update value of single key",
			args: &args{
				mapValue: entity.Map{},
				value:    "test",
				keys: []string{
					"key_a",
				},
			},
			want: &want{
				mapValue: entity.Map{
					"key_a": "test",
				},
				updated: true,
			},
		},
		{
			name: "should not update when intermediary key is not a map",
			args: &args{
				mapValue: entity.Map{
					"key_a": entity.Map{
						"key_b": "tt",
					},
				},
				value: "test",
				keys: []string{
					"key_a",
					"key_b",
					"key_c",
				},
			},
			want: &want{
				mapValue: entity.Map{
					"key_a": entity.Map{
						"key_b": "tt",
					},
				},
				updated: false,
			},
		},
		{
			name: "should update value when intermediary key does not exist",
			args: &args{
				mapValue: entity.Map{},
				value:    "test",
				keys: []string{
					"key_a",
					"key_b",
				},
			},
			want: &want{
				mapValue: entity.Map{
					"key_a": entity.Map{
						"key_b": "test",
					},
				},
				updated: true,
			},
		},
		{
			name: "should update value when keys match structure of map[string]interface{}",
			args: &args{
				mapValue: entity.Map{
					"key_a": map[string]interface{}{},
				},
				value: "test",
				keys: []string{
					"key_a",
					"key_b",
				},
			},
			want: &want{
				mapValue: entity.Map{
					"key_a": entity.Map{
						"key_b": "test",
					},
				},
				updated: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := tt.args.mapValue.Set(
				tt.args.value,
				tt.args.keys...,
			)

			assert.Equal(t, tt.want.updated, updated)
			assert.Equal(t, tt.want.mapValue, tt.args.mapValue)
		})
	}
}

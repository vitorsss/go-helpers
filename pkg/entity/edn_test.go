package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
	"github.com/vitorsss/go-helpers/pkg/entity"
	"olympos.io/encoding/edn"
)

func Test_Set_MarshalEDN(t *testing.T) {
	type args struct {
		set entity.Set[interface{}]
	}
	type want struct {
		result string
		err    error
	}

	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "should marshal simple strings successfully",
			args: &args{
				set: entity.Set[interface{}]{
					"a", "b", "c",
				},
			},
			want: &want{
				result: `#{"a""b""c"}`,
			},
		},
		{
			name: "should marshal simple structs successfully",
			args: &args{
				set: entity.Set[interface{}]{
					testEntity{ID: 1},
					testEntity{ID: 2},
				},
			},
			want: &want{
				result: `#{{:id 1}{:id 2}}`,
			},
		},
		{
			name: "should marshal keywords successfully",
			args: &args{
				set: entity.Set[interface{}]{
					edn.Keyword("test"),
					edn.Keyword("name/spaced"),
				},
			},
			want: &want{
				result: `#{:test :name/spaced}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := edn.Marshal(tt.args.set)
			if assertutil.Error(t, tt.want.err, err) {
				assert.Equal(t, tt.want.result, string(result))
			}
		})
	}
}

func Test_Set_UnmarshalEDN(t *testing.T) {
	t.Run("should unmarshal simple strings", func(t *testing.T) {
		input := entity.Set[string]{}
		err := edn.Unmarshal([]byte(`#{"a""b""c"}`), &input)
		if assertutil.Error(t, nil, err) {
			assert.Equal(t, entity.Set[string]{"a", "b", "c"}, input)
		}
	})
	t.Run("should unmarshal simple structs", func(t *testing.T) {
		input := entity.Set[testEntity]{}
		err := edn.Unmarshal([]byte(`#{{:id 1}{:id 2}}`), &input)
		if assertutil.Error(t, nil, err) {
			assert.Equal(t, entity.Set[testEntity]{
				testEntity{ID: 1},
				testEntity{ID: 2},
			}, input)
		}
	})
}

package edntostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
)

func Test_parseEDNToGolangStructs(t *testing.T) {
	type args struct {
		packagePath string
		prefix      string
		ednContent  []byte
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
			name: "should parse simple namespaced edn",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "SimpleNamespaced",
				ednContent:  []byte(`{:namespace/id "13", :namespace/amount 23, :namespace/sub {:sub/id 1}}`),
			},
			want: &want{
				result: `package path

type SimpleNamespacedNamespace = struct {
	Amount int64                        "json:\"amount\" edn:\"namespace/amount\""
	ID     string                       "json:\"id\" edn:\"namespace/id\""
	Sub    SimpleNamespacedNamespaceSub "json:\"sub\" edn:\"namespace/sub\""
}

type SimpleNamespacedNamespaceSub = struct {
	ID int64 "json:\"id\" edn:\"sub/id\""
}
`,
			},
		},
		{
			name: "should parse default tag values edn",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "DefaultTags",
				ednContent:  []byte(`{:time #inst "2024-05-05T12:29:17Z" :uuid #uuid "4f25e20d-522d-4963-897f-eb04d6d133a2"}`),
			},
			want: &want{
				result: `package path

import (
	time "time"

	uuid "github.com/google/uuid"
)

type DefaultTags = struct {
	Time time.Time "json:\"time\" edn:\"time\""
	Uuid uuid.UUID "json:\"uuid\" edn:\"uuid\""
}
`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEDNToGolangStructs(
				tt.args.packagePath,
				tt.args.prefix,
				tt.args.ednContent,
			)

			if assertutil.Error(t, tt.want.err, err) {
				assert.Equal(t, tt.want.result, string(result))
			}
		})
	}
}

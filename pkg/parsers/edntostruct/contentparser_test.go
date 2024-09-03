package edntostruct

import (
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
	"github.com/vitorsss/go-helpers/pkg/logs"
)

func Test_ContentParser_parseEDNToGolangStructs(t *testing.T) {
	type args struct {
		options     []Option
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
				ednContent:  []byte(`{:namespace/id "13", :namespace/amount 23.19, :namespace/cents 2319, :namespace/valid? false, :namespace/sub {:sub/id 1}}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

type SimpleNamespacedNamespace struct {
	Amount float64                      "json:\"amount\" edn:\"namespace/amount\""
	Cents  int64                        "json:\"cents\" edn:\"namespace/cents\""
	ID     string                       "json:\"id\" edn:\"namespace/id\""
	Sub    SimpleNamespacedNamespaceSub "json:\"sub\" edn:\"namespace/sub\""
	Valid  bool                         "json:\"valid?\" edn:\"namespace/valid?\""
}

type SimpleNamespacedNamespaceSub struct {
	ID int64 "json:\"id\" edn:\"sub/id\""
}
`,
			},
		},
		{
			name: "should parse simple vector and lists",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "Slices",
				ednContent:  []byte(`{:a [{:t 1}], :b [:j], :c ({:h "tt"}), :d (1)}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	edn "olympos.io/encoding/edn"
)

type Slices struct {
	A []SlicesA     "json:\"a\" edn:\"a\""
	B []SlicesBCode "json:\"b\" edn:\"b\""
	C []SlicesC     "json:\"c\" edn:\"c\""
	D []int64       "json:\"d\" edn:\"d\""
}

type SlicesA struct {
	T int64 "json:\"t\" edn:\"t\""
}

type SlicesBCode string

const (
	SlicesBCodeJ SlicesBCode = "j"
)

func (e *SlicesBCode) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}
	*e = SlicesBCode(keyword)
	return err
}

func (e SlicesBCode) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(e))
}

type SlicesC struct {
	H string "json:\"h\" edn:\"h\""
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
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	time "time"

	uuid "github.com/google/uuid"
)

type DefaultTags struct {
	Time time.Time "json:\"time\" edn:\"time\""
	Uuid uuid.UUID "json:\"uuid\" edn:\"uuid\""
}
`,
			},
		}, {
			name: "should parse keyword edn values as 'enums'",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "KeywordEnum",
				ednContent:  []byte(`{:namesss/enum1 :namesss.test/eee :namesss/enum2 :jjj}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	errors "errors"
	fmt "fmt"
	strings "strings"

	edn "olympos.io/encoding/edn"
)

type KeywordEnumNamesss struct {
	Enum1 KeywordEnumNamesssEnum1Code "json:\"enum_1\" edn:\"namesss/enum1\""
	Enum2 KeywordEnumNamesssEnum2Code "json:\"enum_2\" edn:\"namesss/enum2\""
}

type KeywordEnumNamesssEnum1Code string

const (
	KeywordEnumNamesssEnum1CodeEee KeywordEnumNamesssEnum1Code = "eee"
)

func (e *KeywordEnumNamesssEnum1Code) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}
	raw, found := strings.CutPrefix(string(keyword), "namesss.test/")
	if !found {
		return errors.New("KeywordEnumNamesssEnum1Code.UnmarshalEDN: invalid keyword")
	}
	*e = KeywordEnumNamesssEnum1Code(raw)
	return err
}

func (e KeywordEnumNamesssEnum1Code) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(fmt.Sprintf("namesss.test/%s", e)))
}

type KeywordEnumNamesssEnum2Code string

const (
	KeywordEnumNamesssEnum2CodeJjj KeywordEnumNamesssEnum2Code = "jjj"
)

func (e *KeywordEnumNamesssEnum2Code) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}
	*e = KeywordEnumNamesssEnum2Code(keyword)
	return err
}

func (e KeywordEnumNamesssEnum2Code) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(e))
}
`,
			},
		},
		{
			name: "should parse set edn types as 'entity.Set'",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "SetsSetsSets",
				ednContent:  []byte(`{:superset/ok #{{:subtype/id 2 :subtype/name "ttt"}} :superset/string #{"some_string"} :superset/keyword #{:some_keyword} :superset/int #{11 33} :superset/set-of-set #{#{"why" "just" "why?"}}}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	entity "github.com/vitorsss/go-helpers/pkg/entity"
	edn "olympos.io/encoding/edn"
)

type SetsSetsSetsSuperset struct {
	Int      entity.Set[int64]                           "json:\"int\" edn:\"superset/int\""
	Keyword  entity.Set[SetsSetsSetsSupersetKeywordCode] "json:\"keyword\" edn:\"superset/keyword\""
	Ok       entity.Set[SetsSetsSetsSupersetSubtype]     "json:\"ok\" edn:\"superset/ok\""
	SetOfSet entity.Set[entity.Set[string]]              "json:\"set_of_set\" edn:\"superset/set-of-set\""
	String   entity.Set[string]                          "json:\"string\" edn:\"superset/string\""
}

type SetsSetsSetsSupersetKeywordCode string

const (
	SetsSetsSetsSupersetKeywordCodeSomeKeyword SetsSetsSetsSupersetKeywordCode = "some_keyword"
)

func (e *SetsSetsSetsSupersetKeywordCode) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}
	*e = SetsSetsSetsSupersetKeywordCode(keyword)
	return err
}

func (e SetsSetsSetsSupersetKeywordCode) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(e))
}

type SetsSetsSetsSupersetSubtype struct {
	ID   int64  "json:\"id\" edn:\"subtype/id\""
	Name string "json:\"name\" edn:\"subtype/name\""
}
`,
			},
		},
		{
			name: "should be able to overwrite default types",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "Overwrite",
				options: []Option{
					WithTagTypeFn("float64", func() (*types.Package, types.Type) {
						decimalPackage := types.NewPackage("github.com/shopspring/decimal", "decimal")
						return decimalPackage, types.NewNamed(
							types.NewTypeName(
								token.NoPos,
								decimalPackage,
								"Decimal",
								nil,
							),
							nil,
							nil,
						)
					}),
				},
				ednContent: []byte(`{:some/field 2M}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	decimal "github.com/shopspring/decimal"
)

type OverwriteSome struct {
	Field decimal.Decimal "json:\"field\" edn:\"some/field\""
}
`,
			},
		},
		{
			name: "should be able to process complex root node",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "ComplexRoot",
				ednContent:  []byte(`({:some/field 2M})`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

type ComplexRootSome struct {
	Field float64 "json:\"field\" edn:\"some/field\""
}
`,
			},
		},
		{
			name: "should be able to overwrite named types",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "OverwriteNamed",
				options: []Option{
					WithNamedTypeFn("OverwriteNamedSomeOther", func() (*types.Package, *types.Named) {
						decimalPackage := types.NewPackage("github.com/some_org/repo/package", "potato")
						return decimalPackage, types.NewNamed(
							types.NewTypeName(
								token.NoPos,
								decimalPackage,
								"MagicValue",
								nil,
							),
							nil,
							nil,
						)
					}),
					WithNamedTypeFn("OverwriteNamedSomeKeywordCode", func() (*types.Package, *types.Named) {
						decimalPackage := types.NewPackage("github.com/some_org/repo/package", "potato")
						return decimalPackage, types.NewNamed(
							types.NewTypeName(
								token.NoPos,
								decimalPackage,
								"MagicEnum",
								nil,
							),
							nil,
							nil,
						)
					}),
				},
				ednContent: []byte(`{:some/field {:other/value 22} :some/keyword :nn/tt}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	potato "github.com/some_org/repo/package"
)

type OverwriteNamedSome struct {
	Field   potato.MagicValue "json:\"field\" edn:\"some/field\""
	Keyword potato.MagicEnum  "json:\"keyword\" edn:\"some/keyword\""
}
`,
			},
		},
		{
			name: "should be able to parse with mixed namespaces",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "Mixed",
				options:     []Option{},
				ednContent:  []byte(`{:some/field 2 :other/keyword 2}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

type MixedGroup struct {
	MixedOther
	MixedSome
}

type MixedOther struct {
	Keyword int64 "json:\"keyword\" edn:\"other/keyword\""
}

type MixedSome struct {
	Field int64 "json:\"field\" edn:\"some/field\""
}
`,
			},
		},
		{
			name: "should be able to parse with mixed struct types",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "Mixed",
				options:     []Option{},
				ednContent:  []byte(`{:some/field {:sub/a 2 :sub/b 2} :some/other {:sub/b 2 :sub/c 2}}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

type MixedSome struct {
	Field MixedSomeSub "json:\"field\" edn:\"some/field\""
	Other MixedSomeSub "json:\"other\" edn:\"some/other\""
}

type MixedSomeSub struct {
	A int64 "json:\"a\" edn:\"sub/a\""
	B int64 "json:\"b\" edn:\"sub/b\""
	C int64 "json:\"c\" edn:\"sub/c\""
}
`,
			},
		},
		{
			name: "should be able to parse with mixed key types",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "Mixed",
				options:     []Option{},
				ednContent:  []byte(`{:some/field {{:k/s 2} {:sub/a 1 :sub/b 2 :sub/d :a} :m {:sub/b 2 :sub/c 3 :sub/d :b}}}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	edn "olympos.io/encoding/edn"
)

type MixedSome struct {
	Field map[interface{}]MixedSomeSub "json:\"field\" edn:\"some/field\""
}

type MixedSomeKeyK struct {
	S int64 "json:\"s\" edn:\"k/s\""
}

type MixedSomeSub struct {
	A int64             "json:\"a\" edn:\"sub/a\""
	B int64             "json:\"b\" edn:\"sub/b\""
	C int64             "json:\"c\" edn:\"sub/c\""
	D MixedSomeSubDCode "json:\"d\" edn:\"sub/d\""
}

type MixedSomeSubDCode string

const (
	MixedSomeSubDCodeA MixedSomeSubDCode = "a"
	MixedSomeSubDCodeB MixedSomeSubDCode = "b"
)

func (e *MixedSomeSubDCode) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}
	*e = MixedSomeSubDCode(keyword)
	return err
}

func (e MixedSomeSubDCode) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(e))
}
`,
			},
		},
		{
			name: "should be able to parse with struct keys without namespace",
			args: &args{
				packagePath: "example.com/some/package/path",
				prefix:      "Mixed",
				options:     []Option{},
				ednContent:  []byte(`{:some/field {{:s 2} {:sub/a 1 :sub/b 2 :sub/d :a} {:s 3} {:sub/b 2 :sub/c 3 :sub/d :b}}}`),
			},
			want: &want{
				result: `// Code generated by endtostruct DO NOT EDIT
package path

import (
	edn "olympos.io/encoding/edn"
)

type MixedSome struct {
	Field map[MixedSomeKey]MixedSomeSub "json:\"field\" edn:\"some/field\""
}

type MixedSomeKey struct {
	S int64 "json:\"s\" edn:\"s\""
}

type MixedSomeSub struct {
	A int64             "json:\"a\" edn:\"sub/a\""
	B int64             "json:\"b\" edn:\"sub/b\""
	C int64             "json:\"c\" edn:\"sub/c\""
	D MixedSomeSubDCode "json:\"d\" edn:\"sub/d\""
}

type MixedSomeSubDCode string

const (
	MixedSomeSubDCodeA MixedSomeSubDCode = "a"
	MixedSomeSubDCodeB MixedSomeSubDCode = "b"
)

func (e *MixedSomeSubDCode) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}
	*e = MixedSomeSubDCode(keyword)
	return err
}

func (e MixedSomeSubDCode) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(e))
}
`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewContentParser(tt.args.options...)
			destPackage := types.NewPackage(tt.args.packagePath, tt.args.packagePath[strings.LastIndex(tt.args.packagePath, "/")+1:])
			result, err := p.ParseEDNContentToGolang(
				destPackage,
				tt.args.prefix,
				tt.args.ednContent,
			)

			if assertutil.Error(t, tt.want.err, err) {
				assert.Equal(t, tt.want.result, string(result))
			} else {
				logs.Logger.Error().Err(err).Msg("")
			}
		})
	}
}

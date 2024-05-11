package edntostruct

import (
	"go/token"
	"go/types"
	"slices"
	"strings"
)

type fieldTagPair struct {
	field *types.Var
	tag   string
}

func createStructOrderedFields(
	destPackage *types.Package,
	name string,
	fieldTagPairs []fieldTagPair,
) *types.Named {
	slices.SortFunc(fieldTagPairs, func(a, b fieldTagPair) int {
		return strings.Compare(a.field.Name(), b.field.Name())
	})
	fields := make([]*types.Var, 0, len(fieldTagPairs))
	tags := make([]string, 0, len(fieldTagPairs))
	for _, pair := range fieldTagPairs {
		fields = append(fields, pair.field)
		tags = append(tags, pair.tag)
	}
	structType := types.NewStruct(fields, tags)
	typeName := types.NewTypeName(
		token.NoPos,
		destPackage,
		name,
		structType,
	)
	object := types.NewNamed(
		typeName,
		structType,
		nil,
	)
	return object
}

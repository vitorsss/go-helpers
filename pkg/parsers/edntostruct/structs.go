package edntostruct

import (
	"fmt"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
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
		nameCompare := strings.Compare(a.field.Name(), b.field.Name())
		if nameCompare == 0 {
			return strings.Compare(a.field.Type().String(), b.field.Type().String())
		}
		return nameCompare
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

func createStructs(
	destPackage *types.Package,
	options *options,
	prefix string,
	byNamespace map[string][]fieldTagPair,
) (types.Type, error) {
	structs := make([]*types.Named, 0, len(byNamespace))
	for namespace, fields := range byNamespace {
		name := fmt.Sprintf("%s%s", prefix, strcase.ToCamel(namespace))
		var object *types.Named
		if fn, ok := options.namedTypes[name]; ok {
			var importPackage *types.Package
			importPackage, object = fn()
			addImportFixName(destPackage, importPackage)
		} else {
			object = createStructOrderedFields(
				destPackage,
				name,
				fields,
			)
			existingObject := destPackage.Scope().Insert(object.Obj())
			if existingObject != nil {
				return nil, errors.New("unsuported mixed types")
			}
		}
		structs = append(structs,
			object,
		)
	}

	var result types.Type
	if len(structs) == 1 {
		result = structs[0]
	} else {
		structsAsFields := []fieldTagPair{}
		for _, s := range structs {
			structsAsFields = append(structsAsFields, fieldTagPair{
				field: types.NewField(
					token.NoPos,
					destPackage,
					s.Obj().Name(),
					s,
					true,
				),
			})
		}

		object := createStructOrderedFields(
			destPackage,
			prefix,
			structsAsFields,
		)
		existingObject := destPackage.Scope().Insert(object.Obj())
		if existingObject != nil {
			return nil, errors.New("unsuported mixed types")
		}
		result = object
	}

	return result, nil
}

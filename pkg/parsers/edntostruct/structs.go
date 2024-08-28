package edntostruct

import (
	"fmt"
	"go/token"
	"go/types"
	"slices"
	"strings"
	_ "unsafe"

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
	var err error
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
			object, err = addStructToPackage(destPackage, object)
			if err != nil {
				return nil, err
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
			fmt.Sprintf("%sGroup", prefix),
			structsAsFields,
		)
		object, err = addStructToPackage(destPackage, object)
		if err != nil {
			return nil, err
		}
		result = object
	}

	return result, nil
}

//go:linkname nastyScopeInsertOverwrite go/types.(*Scope).insert
func nastyScopeInsertOverwrite(*types.Scope, string, types.Object)

func addStructToPackage(
	destPackage *types.Package,
	object *types.Named,
) (*types.Named, error) {
	existingObject := destPackage.Scope().Insert(object.Obj())
	if existingObject != nil {
		switch existing := existingObject.Type().(type) {
		case *types.Struct:
			newStruct, ok := object.Obj().Type().(*types.Struct)
			if !ok {
				return nil, errors.Errorf("unsuported mixed types: *types.Struct # %s", object.Obj().Type().String())
			}
			var err error
			object, err = mergeStructs(
				destPackage,
				object.Obj().Name(),
				existing,
				newStruct,
			)
			if err != nil {
				return nil, err
			}
			nastyScopeInsertOverwrite(destPackage.Scope(), object.Obj().Name(), object.Obj())
		default:
			return nil, errors.Errorf("unsuported mixed types: %s", existing.String())
		}
	}
	return object, nil
}

func mergeStructs(
	destPackage *types.Package,
	name string,
	existingStruct *types.Struct,
	newStruct *types.Struct,
) (*types.Named, error) {
	fieldTagPairs := []fieldTagPair{}
	existingFields := map[string]*types.Var{}

	for i := 0; i < existingStruct.NumFields(); i++ {
		field := existingStruct.Field(i)
		existingFields[field.Name()] = field
		fieldTagPairs = append(fieldTagPairs, fieldTagPair{
			field: field,
			tag:   existingStruct.Tag(i),
		})
	}

	for i := 0; i < newStruct.NumFields(); i++ {
		field := newStruct.Field(i)
		if existingField, ok := existingFields[field.Name()]; ok {
			if field.Type().String() == existingField.Type().String() {
				continue
			} else {
				return nil, errors.Errorf("types missmatch while merging structs: '%s' != '%s'",
					field.Type().String(),
					existingField.Type().String(),
				)
			}
		}
		fieldTagPairs = append(fieldTagPairs, fieldTagPair{
			field: field,
			tag:   newStruct.Tag(i),
		})
	}

	return createStructOrderedFields(
		destPackage,
		name,
		fieldTagPairs,
	), nil
}

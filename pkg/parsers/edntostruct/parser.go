package edntostruct

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"mvdan.cc/gofumpt/format"
	"olympos.io/encoding/edn"
)

func parseEDNToGolangStructs(
	packagePath string,
	prefix string,
	ednContent []byte,
) ([]byte, error) {
	ednMap := map[interface{}]interface{}{}
	err := edn.Unmarshal(ednContent, &ednMap)
	if err != nil {
		return nil, err
	}

	destPackage := types.NewPackage(packagePath, packagePath[strings.LastIndex(packagePath, "/")+1:])

	_, err = parseEDNTypeToGolangStruct(
		destPackage,
		defaultTagTypes,
		prefix,
		ednMap,
	)
	if err != nil {
		return nil, err
	}
	qualifier := func(other *types.Package) string {
		if destPackage == other {
			return ""
		}
		return other.Name()
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("package %s", destPackage.Name()))
	buffer.WriteString("\n\nimport (\n")
	for _, importParckage := range destPackage.Imports() {
		buffer.WriteString(fmt.Sprintf(`%s "%s"`, importParckage.Name(), importParckage.Path()))
		buffer.WriteString("\n")
	}
	buffer.WriteString(")")
	scope := destPackage.Scope()
	for _, name := range scope.Names() {
		buffer.WriteString("\n\n")
		unformatted := types.ObjectString(scope.Lookup(name), qualifier)
		buffer.Write([]byte(unformatted))
	}

	fmt.Println(buffer.String())

	return format.Source(
		buffer.Bytes(),
		format.Options{
			ModulePath: packagePath,
			ExtraRules: true,
		},
	)
}

type fieldTagPair struct {
	field *types.Var
	tag   string
}

func parseEDNTypeToGolangStruct(
	destPackage *types.Package,
	tagTypes map[string]TypeFn,
	prefix string,
	ednType map[interface{}]interface{},
) (types.Type, error) {
	byNamespace := map[string][]fieldTagPair{}
	for iKey, iVal := range ednType {
		var key string
		switch v := iKey.(type) {
		case string:
			key = v
		case edn.Keyword:
			key = string(v)
		default:
			return nil, errors.New("unmapped key type")
		}
		keyParts := strings.Split(key, "/")
		namespace := ""
		if len(keyParts) > 1 {
			namespace = keyParts[0]
		}

		parsedField, tag, err := parseEDNTypeToGolangField(
			destPackage,
			tagTypes,
			fmt.Sprintf("%s%s", prefix, strcase.ToCamel(namespace)),
			key,
			keyParts[len(keyParts)-1],
			iVal,
		)
		if err != nil {
			return nil, err
		}

		byNamespace[namespace] = append(byNamespace[namespace], fieldTagPair{
			field: parsedField,
			tag:   tag,
		})
	}

	structs := make([]*types.Named, 0, len(byNamespace))
	for namespace, fields := range byNamespace {
		object := createStruct(
			destPackage,
			prefix,
			namespace,
			fields,
		)
		structs = append(structs,
			object,
		)
		existingObject := destPackage.Scope().Insert(object.Obj())
		if existingObject != nil {
			return nil, errors.New("unsuported mixed types")
		}
	}

	var result types.Type
	if len(structs) == 1 {
		result = structs[0]
	} else {
		return nil, errors.New("unsuported mixed namespaces")
	}

	return result, nil
}

func createStruct(
	destPackage *types.Package,
	prefix string,
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
		fmt.Sprintf("%s%s", prefix, strcase.ToCamel(name)),
		structType,
	)
	object := types.NewNamed(
		typeName,
		structType,
		nil,
	)
	return object
}

func parseEDNTypeToGolangField(
	destPackage *types.Package,
	tagTypes map[string]TypeFn,
	prefix string,
	key string,
	name string,
	fieldVal interface{},
) (*types.Var, string, error) {
	var fieldType types.Type
	var varType *types.Var
	var err error
	switch v := fieldVal.(type) {
	case string:
		fieldType = types.Typ[types.String]
	case bool:
		fieldType = types.Typ[types.Bool]
	case map[interface{}]interface{}:
		structBase, err := parseEDNTypeToGolangStruct(
			destPackage,
			tagTypes,
			prefix,
			v,
		)
		if err != nil {
			return nil, "", err
		}
		fieldType = structBase
	case []interface{}:
		if len(v) == 0 {
			return nil, "", errors.New("empty list conversion")
		}
		varType, _, err = parseEDNTypeToGolangField(
			destPackage,
			tagTypes,
			prefix,
			key,
			name,
			v[0],
		)
		if err != nil {
			return nil, "", err
		}
		fieldType = types.NewSlice(varType.Type())
	case map[interface{}]bool:
		if len(v) == 0 {
			return nil, "", errors.New("empty set conversion")
		}
		keys := make([]interface{}, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		varType, _, err = parseEDNTypeToGolangField(
			destPackage,
			tagTypes,
			prefix,
			key,
			name,
			keys[0],
		)
		if err != nil {
			return nil, "", err
		}
		fieldType = nil // TODO: should create an Set type to handle this
	case *interface{}:
		varType, _, err = parseEDNTypeToGolangField(
			destPackage,
			tagTypes,
			prefix,
			key,
			name,
			*v,
		)
		if err != nil {
			return nil, "", err
		}
		fieldType = types.NewPointer(varType.Type())
	case float32, float64, *big.Float:
		fieldType = types.Typ[types.Float64]
	case int, int8, int16, int32, int64:
		fieldType = types.Typ[types.Int64]
	case time.Time:
		typeFn, ok := tagTypes["inst"]
		if !ok {
			return nil, "", errors.New("unmapped tagname")
		}
		var importPackage *types.Package
		importPackage, fieldType = typeFn()
		fmt.Println(importPackage, fieldType)
		addImportFixName(destPackage, importPackage)
	case edn.Tag:
		typeFn, ok := tagTypes[v.Tagname]
		if !ok {
			return nil, "", errors.New("unmapped tagname")
		}
		var importPackage *types.Package
		importPackage, fieldType = typeFn()
		fmt.Println(importPackage, fieldType)
		addImportFixName(destPackage, importPackage)
	case edn.Keyword:
		fieldType = nil // TODO: should create an Enum type to handle this
		// keyParts := strings.Split(string(v), "/")
		// keyName := strcase.ToCamel(keyParts[0])
		// namespaces := map[string]bool{
		// 	keyParts[0]: true,
		// }
		// if len(keyParts) == 1 {
		// 	keyName = strcase.ToCamel(name)
		// 	namespaces = map[string]bool{}
		// }
		// result.Object = &object{
		// 	Name:       fmt.Sprintf("%s%sCode", prefix, keyName),
		// 	Namespaces: namespaces,
		// }
	case nil:
		return nil, "", errors.New("nil value")
	default:
		return nil, "", errors.Errorf("unmapped value type: %#v", v)
	}
	nameCamel := strcase.ToCamel(name)
	if nameCamel == "Id" {
		nameCamel = "ID"
	}
	tag := fmt.Sprintf(`json:"%s" edn:"%s"`, strcase.ToSnake(name), key)
	return types.NewVar(token.NoPos, destPackage, nameCamel, fieldType), tag, nil
}

func addImportFixName(
	destPackage *types.Package,
	importPackage *types.Package,
) {
	imports := destPackage.Imports()
	for _, existingImport := range imports {
		if existingImport.Path() == importPackage.Path() {
			importPackage.SetName(existingImport.Name())
			return
		}
		if existingImport.Name() == importPackage.Name() {
			changed := false
			name := regexp.MustCompile("\\d+").
				ReplaceAllStringFunc(
					importPackage.Name(),
					func(s string) string {
						changed = true
						i, _ := strconv.Atoi(s)
						return strconv.Itoa(i + 1)
					},
				)
			if !changed {
				name = fmt.Sprintf("%s1", name)
			}
			importPackage.SetName(name)
		}
	}
	imports = append(imports, importPackage)
	destPackage.SetImports(imports)
}

package edntostruct

import (
	"fmt"
	"go/token"
	"go/types"
	"math/big"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"olympos.io/encoding/edn"
)

type Parser struct {
	options *options
}

func NewParser(opts ...Option) *Parser {
	opt := defaultOptions()

	for _, optFn := range opts {
		optFn(opt)
	}

	return &Parser{
		options: opt,
	}
}

func (p *Parser) ParseEDNToGolang(
	destPackage *types.Package,
	prefix string,
	ednContent []byte,
) ([]byte, error) {
	ednMap := map[interface{}]interface{}{}
	err := edn.Unmarshal(ednContent, &ednMap)
	if err != nil {
		return nil, err
	}

	_, err = p.parseEDNTypeToGolangStruct(
		destPackage,
		prefix,
		ednMap,
	)
	if err != nil {
		return nil, err
	}

	return printPackage(destPackage)
}

func (p *Parser) parseEDNTypeToGolangStruct(
	destPackage *types.Package,
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

		parsedField, tag, err := p.parseEDNTypeToGolangField(
			destPackage,
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
		object := createStructOrderedFields(
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

func (p *Parser) parseEDNTypeToGolangField(
	destPackage *types.Package,
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
		structBase, err := p.parseEDNTypeToGolangStruct(
			destPackage,
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
		varType, _, err = p.parseEDNTypeToGolangField(
			destPackage,
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
		varType, _, err = p.parseEDNTypeToGolangField(
			destPackage,
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
		varType, _, err = p.parseEDNTypeToGolangField(
			destPackage,
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
		typeFn, ok := p.options.tagTypes["inst"]
		if !ok {
			return nil, "", errors.New("unmapped tagname")
		}
		var importPackage *types.Package
		importPackage, fieldType = typeFn()
		addImportFixName(destPackage, importPackage)
	case edn.Tag:
		typeFn, ok := p.options.tagTypes[v.Tagname]
		if !ok {
			return nil, "", errors.New("unmapped tagname")
		}
		var importPackage *types.Package
		importPackage, fieldType = typeFn()
		addImportFixName(destPackage, importPackage)
	case edn.Keyword:
		keyParts := strings.Split(string(v), "/")
		keyName := name
		namespace := ""
		if len(keyParts) > 1 {
			namespace = keyParts[0]
		}
		fieldType, err = newEnumType(
			destPackage,
			prefix,
			namespace,
			keyName,
			keyParts[len(keyParts)-1],
		)
		if err != nil {
			return nil, "", err
		}
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

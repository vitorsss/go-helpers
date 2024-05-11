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
		name := keyParts[len(keyParts)-1]
		namespace := ""
		if len(keyParts) > 1 {
			namespace = keyParts[0]
		}

		parsedField, tag, err := p.parseEDNTypeToGolangField(
			destPackage,
			fmt.Sprintf("%s%s", prefix, strcase.ToCamel(namespace)),
			namespace,
			name,
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
	namespace string,
	name string,
	fieldVal interface{},
) (*types.Var, string, error) {
	var fieldType types.Type
	var tagType string
	var varType *types.Var
	var err error
	switch v := fieldVal.(type) {
	case map[interface{}]interface{}:
		if namespace == "" {
			prefix = fmt.Sprintf("%s%s", prefix, strcase.ToCamel(name))
		}
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
			namespace,
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
			namespace,
			name,
			keys[0],
		)
		if err != nil {
			return nil, "", err
		}
		elem := varType.Type()
		if pointer, ok := elem.(*types.Pointer); ok {
			elem = pointer.Elem()
		}
		fieldType, err = newSetType(destPackage, elem)
		if err != nil {
			return nil, "", err
		}
	case *interface{}:
		varType, _, err = p.parseEDNTypeToGolangField(
			destPackage,
			prefix,
			namespace,
			name,
			*v,
		)
		if err != nil {
			return nil, "", err
		}
		fieldType = types.NewPointer(varType.Type())
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
	case time.Time:
		tagType = "inst"
	case edn.Tag:
		tagType = v.Tagname
	case float32, float64, *big.Float:
		tagType = "float64"
	case int, int8, int16, int32, int64:
		tagType = "int64"
	case string:
		tagType = "string"
	case bool:
		tagType = "bool"
	case nil:
		return nil, "", errors.New("nil value")
	default:
		return nil, "", errors.Errorf("unmapped value type: %#v", v)
	}
	if tagType != "" {
		typeFn, ok := p.options.tagTypes[tagType]
		if !ok {
			return nil, "", errors.New("unmapped tagname")
		}
		var importPackage *types.Package
		importPackage, fieldType = typeFn()
		if importPackage != nil {
			addImportFixName(destPackage, importPackage)
		}
	}
	nameCamel := strcase.ToCamel(name)
	if nameCamel == "Id" {
		nameCamel = "ID"
	}
	key := name
	if namespace != "" {
		key = fmt.Sprintf("%s/%s", namespace, key)
	}
	tag := fmt.Sprintf(`json:"%s" edn:"%s"`, strcase.ToSnake(name), key)
	return types.NewVar(token.NoPos, destPackage, nameCamel, fieldType), tag, nil
}

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

type ContentParser struct {
	options *options
}

func NewContentParser(opts ...Option) *ContentParser {
	opt := defaultOptions()

	for _, optFn := range opts {
		optFn(opt)
	}

	return &ContentParser{
		options: opt,
	}
}

func (p *ContentParser) ParseEDNContentToGolang(
	destPackage *types.Package,
	prefix string,
	ednContent []byte,
) ([]byte, error) {
	var ednMap interface{}
	err := edn.Unmarshal(ednContent, &ednMap)
	if err != nil {
		return nil, err
	}

	return p.ParseLoadedEDNToGolang(
		destPackage,
		prefix,
		ednMap,
	)
}

func (p *ContentParser) ParseLoadedEDNToGolang(
	destPackage *types.Package,
	prefix string,
	ednMap interface{},
) ([]byte, error) {
	_, _, err := p.parseEDNTypeToGolangField(
		destPackage,
		prefix,
		"",
		"",
		ednMap,
	)
	if err != nil {
		return nil, err
	}

	return printPackage(destPackage)
}

func (p *ContentParser) parseEDNTypeToGolangStruct(
	destPackage *types.Package,
	prefix string,
	parentNamespace string,
	ednType map[interface{}]interface{},
) (types.Type, error) {
	byNamespace := map[string][]fieldTagPair{}
	keyValues := map[string][]interface{}{}
	hasStructKey := false
	keyTypes := []types.Type{}
	for iKey, iVal := range ednType {
		key, keyType, err := p.parseKey(
			destPackage,
			prefix,
			parentNamespace,
			iKey,
		)
		if err != nil {
			return nil, err
		}
		if key == "" {
			hasStructKey = true
		}

		keyTypes = append(keyTypes, keyType)
		keyValues[key] = append(keyValues[key], iVal)
	}

	for key, values := range keyValues {
		for _, iVal := range values {
			keyParts := strings.Split(key, "/")
			name := ""
			namespace := ""
			if !hasStructKey {
				name = keyParts[len(keyParts)-1]
				if len(keyParts) > 1 {
					namespace = keyParts[0]
				}
			}
			sufix := namespace
			if !hasStructKey && sufix == "" && parentNamespace != "" {
				sufix = "unnamespaced"
			}

			parsedField, tag, err := p.parseEDNTypeToGolangField(
				destPackage,
				fmt.Sprintf("%s%s", prefix, strcase.ToCamel(sufix)),
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
	}

	if hasStructKey {
		return createMapType(
			keyTypes,
			byNamespace,
		)
	}
	return createStructs(
		destPackage,
		p.options,
		prefix,
		parentNamespace,
		byNamespace,
	)
}

func (p *ContentParser) parseKey(
	destPackage *types.Package,
	prefix string,
	parentNamespace string,
	iKey interface{},
) (string, types.Type, error) {
	var key string
	var keyType types.Type
	switch v := iKey.(type) {
	case string:
		key = v
		keyType = types.Typ[types.String]
	case edn.Keyword:
		key = string(v)
		keyType = types.Typ[types.String]
	case map[interface{}]interface{}:
		keyType, err := p.parseEDNTypeToGolangStruct(
			destPackage,
			fmt.Sprintf("%sKey", prefix),
			parentNamespace,
			v,
		)
		return "", keyType, err
	case *interface{}:
		return p.parseKey(
			destPackage,
			prefix,
			parentNamespace,
			*v,
		)
	default:
		return "", nil, errors.New("unmapped key type")
	}
	return key, keyType, nil
}

func (p *ContentParser) parseEDNTypeToGolangField(
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
			namespace,
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
		name := fmt.Sprintf("%s%sCode", prefix, strcase.ToCamel(keyName))
		if fn, ok := p.options.namedTypes[name]; ok {
			var importPackage *types.Package
			importPackage, fieldType = fn()
			addImportFixName(destPackage, importPackage)
		} else {
			fieldType, err = newEnumType(
				destPackage,
				namespace,
				name,
				keyParts[len(keyParts)-1],
			)
			if err != nil {
				return nil, "", err
			}
		}
	case edn.Symbol:
		switch v {
		case "Any":
			fieldType = types.NewInterfaceType(nil, nil)
		default:
			return nil, "", errors.New("unmapped symbol type")
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
		addImportFixName(destPackage, importPackage)
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

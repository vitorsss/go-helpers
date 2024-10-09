package edntostruct

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"olympos.io/encoding/edn"
)

type SchemaParser struct {
	options *options
}

func NewSchemaParser(opts ...Option) *SchemaParser {
	opt := defaultOptions()

	for _, optFn := range opts {
		optFn(opt)
	}

	return &SchemaParser{
		options: opt,
	}
}

func (p *SchemaParser) ParseSchemaSchemaToGolang(
	destPackage *types.Package,
	prefix string,
	schemaSchema []byte,
) ([]byte, error) {
	var schemaMap interface{}
	err := edn.Unmarshal(schemaSchema, &schemaMap)
	if err != nil {
		return nil, err
	}

	return p.ParseLoadedSchemaToGolang(
		destPackage,
		prefix,
		schemaMap,
	)
}

func (p *SchemaParser) ParseLoadedSchemaToGolang(
	destPackage *types.Package,
	prefix string,
	schemaMap interface{},
) ([]byte, error) {
	_, _, err := p.parseEDNTypeToGolangField(
		destPackage,
		prefix,
		"",
		"",
		schemaMap,
	)
	if err != nil {
		return nil, err
	}

	return printPackage(destPackage)
}

func (p *SchemaParser) parseEDNTypeToGolangStruct(
	destPackage *types.Package,
	prefix string,
	parentNamespace string,
	schemaType map[interface{}]interface{},
) (types.Type, error) {
	byNamespace := map[string][]fieldTagPair{}
	keyValues := map[string][]interface{}{}
	hasStructKey := false
	keyTypes := []types.Type{}
	for iKey, iVal := range schemaType {
		key, keyType, iVal, err := p.parseKey(
			destPackage,
			prefix,
			parentNamespace,
			iKey,
			iVal,
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
			name := keyParts[len(keyParts)-1]
			namespace := ""
			if len(keyParts) > 1 {
				namespace = keyParts[0]
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

func (p *SchemaParser) parseKey(
	destPackage *types.Package,
	prefix string,
	parentNamespace string,
	iKey interface{},
	iVal interface{},
) (string, types.Type, interface{}, error) {
	var err error
	var key string
	var keyType types.Type
	valResult := iVal
	switch v := iKey.(type) {
	case string:
		key = v
	case edn.Keyword:
		key = string(v)
	case []interface{}:
		switch first := v[0].(type) {
		case edn.Symbol:
			switch first {
			case "optional-key":
				key, keyType, _, err = p.parseKey(
					destPackage,
					prefix,
					parentNamespace,
					v[1],
					iVal,
				)
				valResult = &iVal
			default:
				return "", nil, nil, errors.New("unmapped key modifier symbol")
			}
		default:
			return "", nil, nil, errors.New("unmapped key array first type")
		}
	case map[interface{}]interface{}:
		keyType, err := p.parseEDNTypeToGolangStruct(
			destPackage,
			fmt.Sprintf("%sKey", prefix),
			parentNamespace,
			v,
		)
		return "", keyType, iVal, err
	case *interface{}:
		key, keyType, valResult, err = p.parseKey(
			destPackage,
			prefix,
			parentNamespace,
			*v,
			iVal,
		)
	default:
		return "", nil, nil, errors.New("unmapped key type")
	}
	return key, keyType, valResult, err
}

func (p *SchemaParser) parseEDNTypeToGolangField(
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
		switch len(v) {
		case 0:
			return nil, "", errors.New("empty list conversion")
		case 1:
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
		default:
			first := v[0]
			rest := v[1:]
			switch firstV := first.(type) {
			case edn.Symbol:
				switch firstV {
				case "constrained":
					varType, _, err = p.parseEDNTypeToGolangField(
						destPackage,
						prefix,
						namespace,
						name,
						v[1],
					)
					if err != nil {
						return nil, "", err
					}
					fieldType = varType.Type()
				case "enum":
					restV := []string{}
					namespace = ""
					for _, i := range rest {
						var value string
						switch iv := i.(type) {
						case string:
							value = iv
						case edn.Keyword:
							value = string(iv)
						default:
							return nil, "", errors.New("unmapped enum value type")
						}
						valueParts := strings.Split(value, "/")
						if len(valueParts) > 1 {
							value = valueParts[1]
							if namespace == "" {
								namespace = valueParts[0]
							} else if namespace != valueParts[0] {
								return nil, "", errors.New("mixed enum namespace")
							}
						}
						restV = append(restV, value)
					}
					name := fmt.Sprintf("%s%sCode", prefix, strcase.ToCamel(name))
					if fn, ok := p.options.namedTypes[name]; ok {
						var importPackage *types.Package
						importPackage, fieldType = fn()
						addImportFixName(destPackage, importPackage)
					} else {
						fieldType, err = newEnumType(
							destPackage,
							namespace,
							name,
							restV...,
						)
						if err != nil {
							return nil, "", err
						}
					}
				default:
					return nil, "", errors.New("unmapped modifier symbol")
				}
			default:
				return nil, "", errors.New("unmapped modifier")
			}
		}
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
	case edn.Symbol:
		switch v {
		case "Int":
			tagType = "int64"
		case "java.math.BigDecimal":
			tagType = "float64"
		case "java.time.LocalDateTime":
			tagType = "inst"
		case "Bool":
			tagType = "bool"
		case "Str":
			tagType = "string"
		case "Uuid":
			tagType = "uuid"
		default:
			return nil, "", errors.Errorf("unmapped symbol type: %#v", v)
		}
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

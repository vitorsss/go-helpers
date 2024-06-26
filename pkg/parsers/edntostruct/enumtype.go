package edntostruct

import (
	"bytes"
	"go/token"
	"go/types"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

type enumType struct {
	namespace string
	name      string
	values    []string
}

func (e *enumType) String() string {
	return "string"
}

func (e *enumType) Underlying() types.Type {
	return types.Typ[types.String]
}

func (e *enumType) ExtraString() string {
	buffer := bytes.NewBufferString("\nconst (\n")
	for _, value := range e.values {
		buffer.WriteString("\t")
		buffer.WriteString(e.name)
		buffer.WriteString(strcase.ToCamel(value))
		buffer.WriteString(" = \"")
		buffer.WriteString(value)
		buffer.WriteString("\"\n")
	}
	buffer.WriteString(")\n")
	buffer.WriteString("\nfunc (e *")
	buffer.WriteString(e.name)
	buffer.WriteString(`) UnmarshalEDN(data []byte) error {
	var keyword edn.Keyword
	err := edn.Unmarshal(data, &keyword)
	if err != nil {
		return err
	}`)
	if e.namespace == "" {
		buffer.WriteString("\n\t*e = ")
		buffer.WriteString(e.name)
		buffer.WriteString("(keyword)\n")
	} else {
		buffer.WriteString(`
	raw, found := strings.CutPrefix(string(keyword), "`)
		buffer.WriteString(e.namespace)
		buffer.WriteString(`/")
	if !found {
		return errors.New("`)
		buffer.WriteString(e.name)
		buffer.WriteString(`.UnmarshalEDN: invalid keyword")
	}`)
		buffer.WriteString("\n\t*e = ")
		buffer.WriteString(e.name)
		buffer.WriteString("(raw)\n")
	}

	buffer.WriteString(`	return err
}

func (e `)
	buffer.WriteString(e.name)
	buffer.WriteString(`) MarshalEDN() ([]byte, error) {
	return edn.Marshal(edn.Keyword(`)
	if e.namespace == "" {
		buffer.WriteString("e")
	} else {
		buffer.WriteString(`fmt.Sprintf("`)
		buffer.WriteString(e.namespace)
		buffer.WriteString(`/%s", e)`)
	}
	buffer.WriteString(`))
}`)
	return buffer.String()
}

func newEnumType(
	destPackage *types.Package,
	namespace string,
	name string,
	values ...string,
) (types.Type, error) {
	enumType := &enumType{
		namespace: namespace,
		name:      name,
		values:    values,
	}

	typeName := types.NewTypeName(
		token.NoPos,
		destPackage,
		enumType.name,
		enumType,
	)

	object := types.NewNamed(
		typeName,
		enumType,
		nil,
	)

	existingObject := destPackage.Scope().Insert(object.Obj())
	if existingObject != nil {
		return nil, errors.New("unsuported mixed types")
	}
	addImportFixName(destPackage, EDNPackage)
	if namespace != "" {
		addImportFixName(destPackage, ErrorsPackage)
	}

	return object, nil
}

package edntostruct

import (
	"go/token"
	"go/types"

	"github.com/pkg/errors"
)

var setTypeName = types.NewTypeName(
	token.NoPos,
	EntityPackage,
	"Set",
	nil,
)

func newSetType(
	destPackage *types.Package,
	elem types.Type,
) (types.Type, error) {
	var typeName *types.TypeName
	switch v := elem.(type) {
	case *types.Named:
		typeName = v.Obj()
	default:
		return nil, errors.New("unmapped set type")
	}
	typeParam := types.NewTypeParam(
		typeName,
		nil,
	)
	named := types.NewNamed(
		setTypeName,
		types.NewSlice(elem),
		nil,
	)
	named.SetTypeParams([]*types.TypeParam{
		typeParam,
	})
	addImportFixName(destPackage, EntityPackage)
	return named, nil
}

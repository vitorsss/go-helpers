package edntostruct

import (
	"go/token"
	"go/types"
)

func newSetType(
	destPackage *types.Package,
	elem types.Type,
) (types.Type, error) {
	typeParam := types.NewTypeParam(
		types.NewTypeName(token.NoPos, EntityPackage, "", nil),
		types.NewInterfaceType(nil, nil),
	)
	typeParam.SetConstraint(elem)
	named := types.NewNamed(
		types.NewTypeName(
			token.NoPos,
			EntityPackage,
			"Set",
			nil,
		),
		types.NewSlice(elem),
		nil,
	)
	named.SetTypeParams([]*types.TypeParam{
		typeParam,
	})
	addImportFixName(destPackage, EntityPackage)
	return named, nil
}

package edntostruct

import (
	"go/token"
	"go/types"
)

type TypeFn func() (*types.Package, *types.Named)

func TimeTypeFn() (*types.Package, *types.Named) {
	timePackage := types.NewPackage("time", "time")
	timeTimeType := types.NewNamed(
		types.NewTypeName(
			token.NoPos,
			timePackage,
			"Time",
			nil,
		),
		nil,
		nil,
	)
	return timePackage, timeTimeType
}

func UUIDTypeFn() (*types.Package, *types.Named) {
	uuidPackage := types.NewPackage("github.com/google/uuid", "uuid")
	uuidUUIDType := types.NewNamed(
		types.NewTypeName(
			token.NoPos,
			uuidPackage,
			"UUID",
			nil,
		),
		nil,
		nil,
	)
	return uuidPackage, uuidUUIDType
}

var defaultTagTypes = map[string]TypeFn{
	"inst": TimeTypeFn,
	"uuid": UUIDTypeFn,
}

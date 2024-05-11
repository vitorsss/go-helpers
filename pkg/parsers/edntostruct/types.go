package edntostruct

import (
	"go/token"
	"go/types"
)

type TypeFn func() (*types.Package, types.Type)
type NamedTypeFn func() (*types.Package, *types.Named)

type TypeExtraStringer interface {
	ExtraString() string
}

func TimeTypeFn() (*types.Package, types.Type) {
	return TimePackage, TimeType
}

func UUIDTypeFn() (*types.Package, types.Type) {
	return UUIDPackage, UUIDType
}

func genericTypeFn(t types.Type) TypeFn {
	return func() (*types.Package, types.Type) {
		return nil, t
	}
}

var (
	EDNPackage    = types.NewPackage("olympos.io/encoding/edn", "edn")
	EntityPackage = types.NewPackage("github.com/vitorsss/go-helpers/pkg/entity", "entity")
	ErrorsPackage = types.NewPackage("errors", "errors")
	TimePackage   = types.NewPackage("time", "time")
	UUIDPackage   = types.NewPackage("github.com/google/uuid", "uuid")

	UUIDType = types.NewNamed(
		types.NewTypeName(
			token.NoPos,
			UUIDPackage,
			"UUID",
			nil,
		),
		nil,
		nil,
	)
	TimeType = types.NewNamed(
		types.NewTypeName(
			token.NoPos,
			TimePackage,
			"Time",
			nil,
		),
		nil,
		nil,
	)
)

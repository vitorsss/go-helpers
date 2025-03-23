package edntostruct

import "go/types"

func createMapType(
	keyTypes []types.Type,
	byNamespace map[string][]fieldTagPair,
) (types.Type, error) {
	var mapKeyType types.Type
	for _, keyType := range keyTypes {
		if keyType == nil {
			continue
		}
		if mapKeyType == nil {
			mapKeyType = keyType
			continue
		}
		if mapKeyType.String() != keyType.String() {
			mapKeyType = types.NewInterfaceType(nil, nil)
			break
		}
	}
	var mapValueType types.Type
	for _, namespaceFields := range byNamespace {
		for _, pair := range namespaceFields {
			if mapValueType == nil {
				mapValueType = pair.field.Type()
				continue
			}
			if mapValueType.String() != pair.field.Type().String() {
				mapValueType = types.NewInterfaceType(nil, nil)
			}
		}
	}
	return types.NewMap(mapKeyType, mapValueType), nil
}

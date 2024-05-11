package edntostruct

import (
	"fmt"
	"go/types"
	"strconv"
)

func addImportFixName(
	destPackage *types.Package,
	importPackage *types.Package,
) {
	if destPackage.Path() == importPackage.Path() {
		return
	}
	imports := destPackage.Imports()
	for _, existingImport := range imports {
		if existingImport.Path() == importPackage.Path() {
			importPackage.SetName(existingImport.Name())
			return
		}
		if existingImport.Name() == importPackage.Name() {
			changed := false
			name := packageNameIndexRegex.
				ReplaceAllStringFunc(
					importPackage.Name(),
					func(s string) string {
						changed = true
						i, _ := strconv.Atoi(s)
						return strconv.Itoa(i + 1)
					},
				)
			if !changed {
				name = fmt.Sprintf("%s1", name)
			}
			importPackage.SetName(name)
		}
	}
	imports = append(imports, importPackage)
	destPackage.SetImports(imports)
}

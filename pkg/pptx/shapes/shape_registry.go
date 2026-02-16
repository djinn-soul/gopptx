package shapes

import (
	"strings"
	"sync"
)

// shapeTypeRegistry holds all valid OOXML preset shape type names.
var shapeTypeRegistry = map[string]bool{}

// shapeAliasRegistry maps common aliases to canonical OOXML preset names.
var shapeAliasRegistry = map[string]string{}

var registryOnce sync.Once

// ensureRegistryInitialized populates the registries if they are not already.
func ensureRegistryInitialized() {
	registryOnce.Do(func() {
		initActionShapes()
		initArrowShapes()
		initBasicShapes()
		initCalloutShapes()
		initFlowchartShapes()
		initMathShapes()
		initRectVariantShapes()
		initStarShapes()
	})
}

// registerShapeType adds a shape type to the valid registry.
func registerShapeType(name string) {
	shapeTypeRegistry[name] = true
	lower := strings.ToLower(name)
	if lower != name {
		shapeAliasRegistry[lower] = name
	}
}

// registerShapeAlias maps an alias to a canonical shape type.
func registerShapeAlias(alias, canonical string) {
	shapeAliasRegistry[alias] = canonical
}

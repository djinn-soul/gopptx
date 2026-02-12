package shapes

import "testing"

func TestShapeRegistryAliasTargetsAreRegistered(t *testing.T) {
	for alias, canonical := range shapeAliasRegistry {
		if !shapeTypeRegistry[canonical] {
			t.Fatalf("alias %q points to unregistered shape type %q", alias, canonical)
		}
	}
}

func TestShapeRegistryCanonicalLowerAliasesPreserved(t *testing.T) {
	for canonical := range shapeTypeRegistry {
		lower := canonical
		if normalized := NormalizeShapeType(lower); normalized != canonical {
			t.Fatalf("canonical shape type %q normalized to %q", canonical, normalized)
		}
	}
}

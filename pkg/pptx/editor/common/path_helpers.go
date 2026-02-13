package common

import (
	"path"
	"path/filepath"
)

// ResolveRelationshipTarget calculates the absolute path of a target based on the source part.
// sourcePart: "ppt/slides/slide1.xml"
// targetRel: "../media/image1.png"
// Result: "ppt/media/image1.png"
func ResolveRelationshipTarget(sourcePart, targetRel string) string {
	// If target is already absolute (no starting dot, no slash), treat as relative to root?
	// OOXML Spec: Targets can be relative to the source part's directory.

	dir := path.Dir(sourcePart)
	// path.Join cleans the path, handling ".." correctly.
	// But we need to use forward slashes explicitly just in case.

	joined := path.Join(dir, targetRel)
	return CanonicalPartPath(joined)
}

// MakeRelativePath calculates the relative path from source to target.
// sourcePart: "ppt/slides/slide1.xml"
// targetPart: "ppt/media/image1.png"
// Result: "../media/image1.png"
func MakeRelativePath(sourcePart, targetPart string) string {
	sourceDir := path.Dir(sourcePart)

	// filepath.Rel uses OS separators (backslash on Windows)
	// We must convert inputs to OS standard if they aren't, but path.Dir returns forward slashes.
	// So let's use filepath.Rel but be careful.

	sDir := filepath.FromSlash(sourceDir)
	tPart := filepath.FromSlash(targetPart)

	rel, err := filepath.Rel(sDir, tPart)
	if err != nil {
		// Fallback to absolute if rel fails (shouldn't happen for internal paths)
		return targetPart
	}
	// Convert back to forward slashes for OOXML
	return filepath.ToSlash(rel)
}

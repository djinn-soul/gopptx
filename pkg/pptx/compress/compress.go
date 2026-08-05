package compress

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

// part is one entry of the source package held in memory.
type part struct {
	Name string
	Data []byte
}

// File compresses the PPTX at inPath and writes the result to outPath.
func File(inPath, outPath string, opts Options) (Result, error) {
	source, err := os.ReadFile(inPath) //nolint:gosec // caller-supplied package path
	if err != nil {
		return Result{}, fmt.Errorf("read %s: %w", inPath, err)
	}
	out, result, err := Bytes(source, opts)
	if err != nil {
		return Result{}, err
	}
	if err = os.WriteFile(outPath, out, 0o600); err != nil {
		return Result{}, fmt.Errorf("write %s: %w", outPath, err)
	}
	return result, nil
}

// Bytes compresses an in-memory PPTX package and returns the new package.
func Bytes(data []byte, opts Options) ([]byte, Result, error) {
	parts, err := readParts(data)
	if err != nil {
		return nil, Result{}, err
	}

	quality := opts.Level.ImageQuality()
	out, result, err := compressOnce(parts, opts, quality)
	if err != nil {
		return nil, Result{}, err
	}
	result.OriginalBytes = int64(len(data))

	// Best-effort target size: keep lowering image quality until it fits.
	for opts.TargetSizeBytes > 0 &&
		result.CompressedBytes > opts.TargetSizeBytes &&
		quality-retryQualityStep >= minRetryQuality {
		quality -= retryQualityStep
		out, result, err = compressOnce(parts, opts, quality)
		if err != nil {
			return nil, Result{}, err
		}
		result.OriginalBytes = int64(len(data))
	}

	return out, result, nil
}

func compressOnce(parts []part, opts Options, quality int) ([]byte, Result, error) {
	result := Result{FinalImageQuality: quality}

	dropped := partsToDrop(parts, opts)
	kept := make([]part, 0, len(parts))
	for _, p := range parts {
		if dropped[p.Name] {
			result.RemovedParts = append(result.RemovedParts, p.Name)
			continue
		}
		kept = append(kept, p)
	}
	sort.Strings(result.RemovedParts)

	for i := range kept {
		switch {
		case isImagePart(kept[i].Name):
			newData, recompressed, resized := optimizeImage(kept[i].Data, kept[i].Name, opts.Level, quality)
			if recompressed {
				result.RecompressedImages++
			}
			if resized {
				result.ResizedImages++
			}
			kept[i].Data = newData
		case opts.OptimizeXML && isXMLPart(kept[i].Name):
			kept[i].Data = minifyXML(kept[i].Data)
		}
	}

	if len(result.RemovedParts) > 0 {
		rewriteReferences(kept, result.RemovedParts)
	}

	out, err := writeParts(kept)
	if err != nil {
		return nil, Result{}, err
	}
	result.CompressedBytes = int64(len(out))
	return out, result, nil
}

func readParts(data []byte) ([]part, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open package: %w", err)
	}
	parts := make([]part, 0, len(zr.File))
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, openErr := f.Open()
		if openErr != nil {
			return nil, fmt.Errorf("open part %s: %w", f.Name, openErr)
		}
		content, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read part %s: %w", f.Name, readErr)
		}
		parts = append(parts, part{Name: f.Name, Data: content})
	}
	return parts, nil
}

func writeParts(parts []part) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, p := range parts {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: p.Name, Method: zip.Deflate})
		if err != nil {
			return nil, fmt.Errorf("create part %s: %w", p.Name, err)
		}
		if _, err = w.Write(p.Data); err != nil {
			return nil, fmt.Errorf("write part %s: %w", p.Name, err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("finalize package: %w", err)
	}
	return buf.Bytes(), nil
}

// partsToDrop decides which parts the options remove.
func partsToDrop(parts []part, opts Options) map[string]bool {
	dropped := make(map[string]bool)
	for _, p := range parts {
		name := normalizeName(p.Name)
		switch {
		case opts.RemoveNotes && strings.HasPrefix(name, "ppt/notesSlides/"):
			dropped[p.Name] = true
		case opts.RemoveComments && isCommentPart(name):
			dropped[p.Name] = true
		case opts.RemoveProperties && isOptionalPropertyPart(name):
			dropped[p.Name] = true
		}
	}
	if opts.RemoveUnusedMedia {
		for name := range unusedMedia(parts, dropped) {
			dropped[name] = true
		}
	}
	// A dropped part's own relationship file goes with it.
	for _, p := range parts {
		if dropped[relsOwner(p.Name)] {
			dropped[p.Name] = true
		}
	}
	return dropped
}

func isCommentPart(name string) bool {
	return strings.HasPrefix(name, "ppt/comments/") ||
		name == "ppt/commentAuthors.xml" ||
		strings.HasPrefix(name, "ppt/modernComments/")
}

func isOptionalPropertyPart(name string) bool {
	return name == "docProps/custom.xml" || strings.HasPrefix(name, "docProps/thumbnail.")
}

// unusedMedia returns media parts that no surviving relationship points at.
func unusedMedia(parts []part, alreadyDropped map[string]bool) map[string]bool {
	referenced := make(map[string]bool)
	for _, p := range parts {
		if !isRelsPart(p.Name) || alreadyDropped[p.Name] {
			continue
		}
		owner := relsOwner(p.Name)
		for _, target := range relationshipTargets(p.Data) {
			if strings.Contains(target, "://") {
				continue
			}
			referenced[resolveTarget(owner, target)] = true
		}
	}

	unused := make(map[string]bool)
	for _, p := range parts {
		name := normalizeName(p.Name)
		if !strings.HasPrefix(name, "ppt/media/") {
			continue
		}
		if !referenced[name] {
			unused[p.Name] = true
		}
	}
	return unused
}

// rewriteReferences removes relationship entries and content-type overrides
// that point at parts which no longer exist.
func rewriteReferences(parts []part, removed []string) {
	removedSet := make(map[string]bool, len(removed))
	for _, name := range removed {
		removedSet[normalizeName(name)] = true
	}
	for i := range parts {
		switch {
		case isRelsPart(parts[i].Name):
			owner := relsOwner(parts[i].Name)
			parts[i].Data = dropRelationships(parts[i].Data, owner, removedSet)
		case normalizeName(parts[i].Name) == "[Content_Types].xml":
			parts[i].Data = dropContentTypeOverrides(parts[i].Data, removedSet)
		}
	}
}

func normalizeName(name string) string {
	return strings.TrimPrefix(strings.ReplaceAll(name, "\\", "/"), "/")
}

func isRelsPart(name string) bool {
	n := normalizeName(name)
	return strings.HasSuffix(n, ".rels") && strings.Contains(n, "_rels/")
}

// relsOwner maps `ppt/slides/_rels/slide1.xml.rels` to `ppt/slides/slide1.xml`.
// It returns an empty string for anything that is not a rels part.
func relsOwner(name string) string {
	n := normalizeName(name)
	if !strings.HasSuffix(n, ".rels") {
		return ""
	}
	dir, file := path.Split(n)
	dir = strings.TrimSuffix(dir, "/")
	if !strings.HasSuffix(dir, "_rels") {
		return ""
	}
	ownerDir := strings.TrimSuffix(strings.TrimSuffix(dir, "_rels"), "/")
	owner := strings.TrimSuffix(file, ".rels")
	if ownerDir == "" {
		return owner
	}
	return ownerDir + "/" + owner
}

// resolveTarget resolves a relationship target against the part that owns the
// relationship, so `../media/image1.png` from `ppt/slides/slide1.xml` becomes
// `ppt/media/image1.png`.
func resolveTarget(owner, target string) string {
	t := strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	if strings.HasPrefix(t, "/") {
		return strings.TrimPrefix(t, "/")
	}
	base := path.Dir(normalizeName(owner))
	if base == "." || base == "" {
		return path.Clean(t)
	}
	return path.Clean(base + "/" + t)
}

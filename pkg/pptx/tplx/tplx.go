package tplx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Context is a map of template variable names to values.
// Values may be strings, bools, ints, float64, or []Row for loop data.
type Context map[string]any

// Row represents a single data record for table-row or slide loops.
type Row map[string]string

// Result holds the rendered PPTX bytes and exposes a Save helper.
type Result struct {
	data []byte
}

// Bytes returns the raw rendered PPTX bytes.
func (r *Result) Bytes() []byte { return r.data }

// Save writes the rendered PPTX to path.
func (r *Result) Save(path string) error {
	return os.WriteFile(path, r.data, 0o600)
}

// Options controls rendering behaviour.
type Options struct {
	// Strict causes Render to return an error when a {{key}} token is not found in ctx.
	// Default: false (lenient — unresolved tokens are left as-is).
	Strict bool
}

var (
	escapedTokenOpen  = []byte("@@TPLX_ESC_OPEN@@")
	escapedTokenClose = []byte("@@TPLX_ESC_CLOSE@@")
)

// Render opens a PPTX template file and renders it with the given context.
func Render(path string, ctx Context) (*Result, error) {
	return RenderWithOptions(path, ctx, Options{})
}

// RenderWithOptions is like Render but accepts rendering options.
func RenderWithOptions(path string, ctx Context, opts Options) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tplx: read template %q: %w", path, err)
	}
	return RenderBytesWithOptions(data, ctx, opts)
}

// RenderBytes renders a PPTX template from in-memory bytes.
func RenderBytes(pptxBytes []byte, ctx Context) (*Result, error) {
	return RenderBytesWithOptions(pptxBytes, ctx, Options{})
}

// RenderBytesWithOptions is like RenderBytes but accepts rendering options.
func RenderBytesWithOptions(pptxBytes []byte, ctx Context, opts Options) (*Result, error) {
	zr, err := zip.NewReader(bytes.NewReader(pptxBytes), int64(len(pptxBytes)))
	if err != nil {
		return nil, fmt.Errorf("tplx: not a valid PPTX/ZIP: %w", err)
	}

	// Build a mutable map of part-name → bytes.
	parts := make(map[string][]byte, len(zr.File))
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err2 := f.Open()
		if err2 != nil {
			return nil, fmt.Errorf("tplx: open %s: %w", f.Name, err2)
		}
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rc)
		_ = rc.Close()
		parts[f.Name] = buf.Bytes()
	}

	// Identify slide parts.
	slideParts := collectSlideParts(parts)
	protectEscapedTemplateDelimiters(parts, slideParts)

	// Pass 1: expand slide-level {{#each}} loops (may add/remove parts).
	slideParts, parts = expandSlideLoops(slideParts, parts, ctx)

	// Pass 2: process each slide — conditionals → table-row loops → scalar replace.
	for _, name := range slideParts {
		data, ok := parts[name]
		if !ok {
			continue
		}
		data = applyConditionals(data, ctx)
		data = expandTableRows(data, ctx)
		data, err = interpolateXMLPart(data, ctx, nil, opts.Strict)
		if err != nil {
			return nil, fmt.Errorf("tplx: interpolate %s: %w", name, err)
		}
		if opts.Strict {
			if tok := firstTemplateToken(data); tok != "" {
				return nil, fmt.Errorf("tplx: unresolved token %q in %s", tok, name)
			}
		}
		parts[name] = data
	}

	// Pass 3: interpolate non-slide text parts (notes, masters, layouts).
	for name, data := range parts {
		if isNonSlideTextPart(name) {
			data, err = interpolateXMLPart(data, ctx, nil, opts.Strict)
			if err != nil {
				return nil, fmt.Errorf("tplx: interpolate %s: %w", name, err)
			}
			if opts.Strict {
				if tok := firstTemplateToken(data); tok != "" {
					return nil, fmt.Errorf("tplx: unresolved token %q in %s", tok, name)
				}
			}
			parts[name] = data
		}
	}

	// Restore escaped literal delimiters on final text parts.
	restoreEscapedTemplateDelimiters(parts, slideParts)

	// Repack into a new ZIP.
	out, err := repackZip(zr, parts)
	if err != nil {
		return nil, fmt.Errorf("tplx: repack: %w", err)
	}
	return &Result{data: out}, nil
}

func firstTemplateToken(data []byte) string {
	match := tokenPattern.Find(data)
	if len(match) == 0 {
		return ""
	}
	return string(match)
}

func protectEscapedTemplateDelimiters(parts map[string][]byte, slideParts []string) {
	for _, name := range slideParts {
		if data, ok := parts[name]; ok {
			parts[name] = protectEscapedTokenText(data)
		}
	}
	for name, data := range parts {
		if isNonSlideTextPart(name) {
			parts[name] = protectEscapedTokenText(data)
		}
	}
}

func restoreEscapedTemplateDelimiters(parts map[string][]byte, slideParts []string) {
	for _, name := range slideParts {
		if data, ok := parts[name]; ok {
			parts[name] = restoreEscapedTokenText(data)
		}
	}
	for name, data := range parts {
		if isNonSlideTextPart(name) {
			parts[name] = restoreEscapedTokenText(data)
		}
	}
}

func protectEscapedTokenText(data []byte) []byte {
	out := bytes.ReplaceAll(data, []byte(`\{{`), escapedTokenOpen)
	out = bytes.ReplaceAll(out, []byte(`\}}`), escapedTokenClose)
	return out
}

func restoreEscapedTokenText(data []byte) []byte {
	out := bytes.ReplaceAll(data, escapedTokenOpen, []byte("{{"))
	out = bytes.ReplaceAll(out, escapedTokenClose, []byte("}}"))
	return out
}

// ── internal helpers ──────────────────────────────────────────────────────────

// collectSlideParts returns all ppt/slides/slideN.xml paths in order.
func collectSlideParts(parts map[string][]byte) []string {
	var slides []string
	for name := range parts {
		if strings.HasPrefix(name, "ppt/slides/slide") && strings.HasSuffix(name, ".xml") &&
			!strings.Contains(name, "_rels") {
			slides = append(slides, name)
		}
	}
	// Sort deterministically.
	sortStrings(slides)
	return slides
}

// expandSlideLoops expands {{#each KEY}} slide-level loops.
// Returns the updated slide list and parts map.
func expandSlideLoops(
	slideParts []string,
	parts map[string][]byte,
	ctx Context,
) ([]string, map[string][]byte) {
	var newSlides []string
	nextNum := maxSlideNumber(slideParts) + 1

	for _, name := range slideParts {
		data := parts[name]
		cond := detectEachSlide(data)
		if cond == nil {
			newSlides = append(newSlides, name)
			continue
		}

		// Get the rows.
		rowsAny, ok := ctx[cond.key]
		if !ok {
			// No data — remove the template slide entirely.
			delete(parts, name)
			deleteSlideRels(parts, name)
			continue
		}
		rows, ok := toRows(rowsAny)
		if !ok {
			delete(parts, name)
			deleteSlideRels(parts, name)
			continue
		}

		expanded := expandSlide(data, rows, ctx)
		// Remove template slide.
		delete(parts, name)
		deleteSlideRels(parts, name)

		// Write one slide part per expanded row.
		for _, slideXML := range expanded {
			newName := fmt.Sprintf("ppt/slides/slide%d.xml", nextNum)
			parts[newName] = slideXML
			newSlides = append(newSlides, newName)
			// Copy the rels file from the original template slide.
			copySlideRels(parts, name, newName)
			nextNum++
		}
	}
	return newSlides, parts
}

func maxSlideNumber(slides []string) int {
	max := 0
	for _, s := range slides {
		n := parseSlideNumber(s)
		if n > max {
			max = n
		}
	}
	return max
}

func parseSlideNumber(name string) int {
	// ppt/slides/slideN.xml
	base := trimPrefix(name, "ppt/slides/slide")
	base = trimSuffix(base, ".xml")
	n := 0
	for _, c := range base {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func deleteSlideRels(parts map[string][]byte, slideName string) {
	relsName := slideRelsPath(slideName)
	delete(parts, relsName)
}

func copySlideRels(parts map[string][]byte, src, dst string) {
	relsData, ok := parts[slideRelsPath(src)]
	if !ok {
		return
	}
	parts[slideRelsPath(dst)] = relsData
}

func slideRelsPath(slideName string) string {
	// ppt/slides/slideN.xml → ppt/slides/_rels/slideN.xml.rels
	base := lastSegment(slideName)
	return "ppt/slides/_rels/" + base + ".rels"
}

func lastSegment(s string) string {
	idx := strings.LastIndex(s, "/")
	if idx < 0 {
		return s
	}
	return s[idx+1:]
}

// isNonSlideTextPart returns true for XML parts that may contain template tokens.
func isNonSlideTextPart(name string) bool {
	if strings.Contains(name, "_rels") {
		return false
	}
	return strings.HasSuffix(name, ".xml") &&
		(strings.HasPrefix(name, "ppt/notesSlides/") ||
			strings.HasPrefix(name, "ppt/slideLayouts/") ||
			name == "ppt/presentation.xml")
}

// repackZip rebuilds the ZIP archive, replacing parts with modified data.
func repackZip(original *zip.Reader, parts map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	written := make(map[string]bool)
	for _, f := range original.File {
		if f.FileInfo().IsDir() {
			continue
		}
		hdr := f.FileHeader
		w, err := zw.CreateHeader(&hdr)
		if err != nil {
			return nil, err //nolint:wrapcheck
		}

		if data, ok := parts[f.Name]; ok {
			if _, err2 := w.Write(data); err2 != nil {
				return nil, err2 //nolint:wrapcheck
			}
		} else {
			rc, err2 := f.Open()
			if err2 != nil {
				return nil, err2 //nolint:wrapcheck
			}
			var tmp bytes.Buffer
			_, _ = tmp.ReadFrom(rc)
			_ = rc.Close()
			if _, err2 = w.Write(tmp.Bytes()); err2 != nil {
				return nil, err2 //nolint:wrapcheck
			}
		}
		written[f.Name] = true
	}

	// Write any new parts not in the original archive (expanded slides).
	for name, data := range parts {
		if written[name] {
			continue
		}
		w, err := zw.Create(name)
		if err != nil {
			return nil, err //nolint:wrapcheck
		}
		if _, err2 := w.Write(data); err2 != nil {
			return nil, err2 //nolint:wrapcheck
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err //nolint:wrapcheck
	}
	return buf.Bytes(), nil
}

// ── tiny stdlib-free sort and string helpers ──────────────────────────────────

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}

func trimPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

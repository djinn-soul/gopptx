package tplx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
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

	parts, err := readZipParts(zr)
	if err != nil {
		return nil, err
	}

	// Identify slide parts.
	slideParts := collectSlideParts(parts)
	protectEscapedTemplateDelimiters(parts, slideParts)

	// Pass 1: expand slide-level {{#each}} loops (may add/remove parts).
	slideParts, parts = expandSlideLoops(slideParts, parts, ctx)

	if err = renderSlideParts(parts, slideParts, ctx, opts.Strict); err != nil {
		return nil, err
	}

	if err = renderNonSlideParts(parts, ctx, opts.Strict); err != nil {
		return nil, err
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

func readZipParts(zr *zip.Reader) (map[string][]byte, error) {
	parts := make(map[string][]byte, len(zr.File))
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		partData, err := readZipFileBytes(f)
		if err != nil {
			return nil, fmt.Errorf("tplx: open %s: %w", f.Name, err)
		}
		parts[f.Name] = partData
	}
	return parts, nil
}

func readZipFileBytes(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err //nolint:wrapcheck // Caller wraps with part-name context.
	}
	defer func() {
		_ = rc.Close()
	}()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(rc)
	return buf.Bytes(), nil
}

func renderSlideParts(parts map[string][]byte, slideParts []string, ctx Context, strict bool) error {
	for _, name := range slideParts {
		data, ok := parts[name]
		if !ok {
			continue
		}
		updated, err := renderSlidePart(data, ctx, strict, name)
		if err != nil {
			return err
		}
		parts[name] = updated
	}
	return nil
}

func renderSlidePart(data []byte, ctx Context, strict bool, name string) ([]byte, error) {
	data = applyConditionals(data, ctx)
	data = expandTableRows(data, ctx)
	data = interpolateXMLPart(data, ctx, nil, strict)
	if err := validateStrictTokens(data, strict, name); err != nil {
		return nil, err
	}
	return data, nil
}

func renderNonSlideParts(parts map[string][]byte, ctx Context, strict bool) error {
	for name, data := range parts {
		if !isNonSlideTextPart(name) {
			continue
		}
		data = interpolateXMLPart(data, ctx, nil, strict)
		if err := validateStrictTokens(data, strict, name); err != nil {
			return err
		}
		parts[name] = data
	}
	return nil
}

func validateStrictTokens(data []byte, strict bool, partName string) error {
	if !strict {
		return nil
	}
	if tok := firstTemplateToken(data); tok != "" {
		return fmt.Errorf("tplx: unresolved token %q in %s", tok, partName)
	}
	return nil
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
	out := bytes.ReplaceAll(data, []byte(`\{{`), escapedTokenOpenBytes())
	out = bytes.ReplaceAll(out, []byte(`\}}`), escapedTokenCloseBytes())
	return out
}

func restoreEscapedTokenText(data []byte) []byte {
	out := bytes.ReplaceAll(data, escapedTokenOpenBytes(), []byte("{{"))
	out = bytes.ReplaceAll(out, escapedTokenCloseBytes(), []byte("}}"))
	return out
}

func escapedTokenOpenBytes() []byte {
	return []byte("@@TPLX_ESC_OPEN@@")
}

func escapedTokenCloseBytes() []byte {
	return []byte("@@TPLX_ESC_CLOSE@@")
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
	maxSlideNum := 0
	for _, s := range slides {
		n := parseSlideNumber(s)
		if n > maxSlideNum {
			maxSlideNum = n
		}
	}
	return maxSlideNum
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
		n = n*numericBaseTen + int(c-'0')
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
		err := writeArchiveEntry(zw, f, parts)
		if err != nil {
			return nil, err
		}
		written[f.Name] = true
	}

	if err := writeNewArchiveParts(zw, parts, written); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err //nolint:wrapcheck // Preserve ZIP finalization errors.
	}
	return buf.Bytes(), nil
}

func writeArchiveEntry(zw *zip.Writer, f *zip.File, parts map[string][]byte) error {
	hdr := f.FileHeader
	w, err := zw.CreateHeader(&hdr)
	if err != nil {
		return err //nolint:wrapcheck // ZIP writer header failures are surfaced directly.
	}
	if data, ok := parts[f.Name]; ok {
		return writeArchiveBytes(w, data)
	}
	data, err := readZipFileBytes(f)
	if err != nil {
		return err //nolint:wrapcheck // Preserve ZIP entry open failures for passthrough files.
	}
	return writeArchiveBytes(w, data)
}

func writeArchiveBytes(w io.Writer, data []byte) error {
	if _, err := w.Write(data); err != nil {
		return err //nolint:wrapcheck // Preserve write failures for generated and passthrough payloads.
	}
	return nil
}

func writeNewArchiveParts(zw *zip.Writer, parts map[string][]byte, written map[string]bool) error {
	for name, data := range parts {
		if written[name] {
			continue
		}
		w, err := zw.Create(name)
		if err != nil {
			return err //nolint:wrapcheck // ZIP writer create failures are surfaced directly.
		}
		if _, err = w.Write(data); err != nil {
			return err //nolint:wrapcheck // Preserve write failures for newly added parts.
		}
	}
	return nil
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

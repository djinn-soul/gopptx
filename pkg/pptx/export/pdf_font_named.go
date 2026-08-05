package export

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/signintech/gopdf"
)

// A deck names its typefaces — "Georgia", "Verdana", "Segoe UI" — and the
// renderer used to collapse every one of them onto three generic faces, so a
// slide set in Georgia came out in whatever the serif fallback happened to be.
// This file indexes the fonts actually installed on the host by their family
// name, so a run can be drawn in the face it asks for and only falls back to the
// generic aliases when the host does not have that family.

// namedFontFaces holds the file for each style of one installed family.
type namedFontFaces struct {
	// paths is keyed by the gopdf style bitmask: Regular, Bold, Italic, or
	// Bold|Italic.
	paths map[int]string
}

//nolint:gochecknoglobals // The installed-font index is host state, shared by every document.
var (
	systemFontIndexMu    sync.Mutex
	systemFontIndex      map[string]*namedFontFaces
	systemFontIndexAge   time.Time
	systemFontIndexBuilt bool
)

// systemFontIndexTTL is how stale the index may get before a miss rebuilds it.
// Rebuilding on every miss would re-scan the font directories for each unknown
// typeface a deck names; never rebuilding would miss a font installed while the
// process was running.
const systemFontIndexTTL = 30 * time.Second

// reservedFontAliases are the aliases the generic faces already occupy. A deck
// that names a typeface "sans" must not be embedded over one of them.
//
//nolint:gochecknoglobals // Immutable set.
var reservedFontAliases = map[string]bool{
	fontFamilySans:  true,
	fontFamilySerif: true,
	fontFamilyMono:  true,
	fontFamilyCJK:   true,
}

// namedFontAliasPrefix disambiguates a family whose name collides with one of
// the generic aliases. Ordinary families are registered under their own name so
// that the PDF's /BaseFont reads "Georgia" rather than an internal token.
const namedFontAliasPrefix = "family-"

// normalizeFontFamily is the key both the index and the alias cache use. PPTX
// typeface names differ from the installed family only in case and spacing.
func normalizeFontFamily(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// cachedNamedPDFFontAlias returns the alias a family was already registered
// under in this document, or "" when it is unregistered or known to be missing.
// It never touches the filesystem, so metrics lookups can use it freely.
func cachedNamedPDFFontAlias(pdf *gopdf.GoPdf, family string) string {
	key := normalizeFontFamily(family)
	if key == "" {
		return ""
	}
	alias, _ := documentFontsFor(pdf).namedAlias(key)
	return alias
}

// ensureNamedPDFFont registers the installed family under its own alias and
// returns that alias, or "" when the host has no such family. Registration is
// lazy: a deck that never mentions Georgia never pays to embed it.
func ensureNamedPDFFont(pdf *gopdf.GoPdf, family string) string {
	key := normalizeFontFamily(family)
	if key == "" {
		return ""
	}
	fonts := documentFontsFor(pdf)
	if alias, ok := fonts.namedAlias(key); ok {
		return alias
	}

	alias := registerNamedPDFFont(pdf, key)
	fonts.putNamedAlias(key, alias)
	return alias
}

// namedFontAlias is the alias one family is embedded under: its own name,
// unless that would collide with a generic alias.
func namedFontAlias(key string) string {
	if reservedFontAliases[key] {
		return namedFontAliasPrefix + key
	}
	return key
}

// registerNamedPDFFont embeds every style of one installed family. The regular
// face is required: gopdf resolves an unstyled SetFont against it, and its
// metrics are the ones the layout code reads.
func registerNamedPDFFont(pdf *gopdf.GoPdf, key string) string {
	faces := lookupSystemFontFaces(key)
	if faces == nil {
		return ""
	}
	regular, ok := faces.paths[gopdf.Regular]
	if !ok {
		return ""
	}
	alias := namedFontAlias(key)
	if err := pdf.AddTTFFontWithOption(alias, regular, nativePDFFontOption(gopdf.Regular)); err != nil {
		return ""
	}
	if err := pdf.SetFont(alias, "", defaultFontSize); err != nil {
		return ""
	}
	registerPDFFontMetrics(pdf, alias, regular)

	for _, style := range []int{gopdf.Bold, gopdf.Italic, gopdf.Bold | gopdf.Italic} {
		path, has := faces.paths[style]
		if !has {
			continue
		}
		if err := pdf.AddTTFFontWithOption(alias, path, nativePDFFontOption(style)); err != nil {
			continue
		}
		// gopdf falls back to the regular face for a style it cannot set, which
		// is the behaviour wanted when a family ships no italic.
		_ = pdf.SetFontWithStyle(alias, style, defaultFontSize)
	}
	return alias
}

// lookupSystemFontFaces finds one installed family. A miss rebuilds the index
// if it has gone stale, so a font installed while the process was running is
// picked up rather than being treated as missing forever.
func lookupSystemFontFaces(key string) *namedFontFaces {
	systemFontIndexMu.Lock()
	defer systemFontIndexMu.Unlock()

	if !systemFontIndexBuilt {
		refreshSystemFontIndexLocked()
	}
	if faces, ok := systemFontIndex[key]; ok {
		return faces
	}
	if time.Since(systemFontIndexAge) >= systemFontIndexTTL {
		refreshSystemFontIndexLocked()
	}
	return systemFontIndex[key]
}

func refreshSystemFontIndexLocked() {
	systemFontIndex = buildSystemFontIndex()
	systemFontIndexAge = time.Now()
	systemFontIndexBuilt = true
}

// buildSystemFontIndex reads the family and style out of every font file in the
// host's font directories. Only each file's name table is read, not the whole
// font, so indexing a few hundred faces stays cheap.
func buildSystemFontIndex() map[string]*namedFontFaces {
	index := make(map[string]*namedFontFaces)
	for _, dir := range systemFontDirs() {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !isFontFileName(entry.Name()) {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			family, style, nameErr := readTTFFamilyAndStyle(path)
			if nameErr != nil || family == "" {
				continue
			}
			addFontToIndex(index, family, style, path)
		}
	}
	return index
}

// addFontToIndex records one face. The first file found for a style wins, so an
// earlier (more specific, per-user) font directory takes precedence.
func addFontToIndex(index map[string]*namedFontFaces, family string, style int, path string) {
	key := normalizeFontFamily(family)
	faces, ok := index[key]
	if !ok {
		faces = &namedFontFaces{paths: map[int]string{}}
		index[key] = faces
	}
	if _, exists := faces.paths[style]; !exists {
		faces.paths[style] = path
	}
}

func isFontFileName(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".ttf", ".otf", ".ttc":
		return true
	default:
		return false
	}
}

// systemFontDirs lists where the host keeps its fonts, most specific first.
func systemFontDirs() []string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		winDir := os.Getenv("WINDIR")
		if winDir == "" {
			winDir = `C:\Windows`
		}
		dirs := []string{}
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			dirs = append(dirs, filepath.Join(local, `Microsoft\Windows\Fonts`))
		}
		return append(dirs, filepath.Join(winDir, "Fonts"))
	case "darwin":
		dirs := []string{}
		if home != "" {
			dirs = append(dirs, filepath.Join(home, "Library", "Fonts"))
		}
		return append(dirs,
			"/Library/Fonts",
			"/System/Library/Fonts",
			"/System/Library/Fonts/Supplemental",
		)
	default:
		dirs := []string{}
		if home != "" {
			dirs = append(dirs, filepath.Join(home, ".fonts"), filepath.Join(home, ".local/share/fonts"))
		}
		return append(dirs,
			"/usr/local/share/fonts",
			"/usr/share/fonts",
			"/usr/share/fonts/truetype",
			"/usr/share/fonts/opentype",
			"/usr/share/fonts/TTF",
		)
	}
}

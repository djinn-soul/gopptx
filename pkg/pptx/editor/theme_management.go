package editor

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	themeOwnerSlideMaster = "slideMaster"
	themeOwnerNotesMaster = "notesMaster"
)

var (
	themePartPathPattern = regexp.MustCompile(`^ppt/theme/theme\d+\.xml$`)
	majorFontPattern     = regexp.MustCompile(`(?is)<(?:\w+:)?majorFont\b[^>]*>.*?</(?:\w+:)?majorFont>`)
	minorFontPattern     = regexp.MustCompile(`(?is)<(?:\w+:)?minorFont\b[^>]*>.*?</(?:\w+:)?minorFont>`)
	latinPattern         = regexp.MustCompile(`(?is)<(?:\w+:)?latin\b[^>]*\/>`)
	typefacePattern      = regexp.MustCompile(`\btypeface="[^"]*"`)
)

// ThemeBinding maps a master owner to a concrete theme part.
type ThemeBinding struct {
	OwnerType      string
	OwnerPart      string
	RelationshipID string
	ThemePart      string
}

// ThemeInventory describes discovered theme parts and owner bindings.
type ThemeInventory struct {
	ThemeParts []string
	Bindings   []ThemeBinding
}

// ThemeColorScheme captures the 12 standard OOXML theme color slots.
type ThemeColorScheme struct {
	Dk1      string
	Lt1      string
	Dk2      string
	Lt2      string
	Accent1  string
	Accent2  string
	Accent3  string
	Accent4  string
	Accent5  string
	Accent6  string
	Hlink    string
	FolHlink string
}

// GetThemeInventory discovers all theme parts and master/theme bindings in the package.
//
//nolint:gocognit // Inventory assembly handles legacy/current theme parts with explicit fallbacks and normalization.
func (e *PresentationEditor) GetThemeInventory() (ThemeInventory, error) {
	if e == nil || e.parts == nil {
		return ThemeInventory{}, errors.New("editor cannot be nil")
	}

	inv := ThemeInventory{ThemeParts: e.parts.KeysWithPrefix("ppt/theme/")}
	sort.Strings(inv.ThemeParts)

	collect := func(prefix, ownerType string) error {
		for _, relPath := range e.parts.KeysWithPrefix(prefix) {
			data, ok := e.parts.Get(relPath)
			if !ok {
				continue
			}
			rels, err := parseRelationshipsXML(data)
			if err != nil {
				return fmt.Errorf("parse %s: %w", relPath, err)
			}
			ownerPart := sourcePartFromRels(relPath)
			for _, rel := range rels {
				if rel.Type != common.RelTypeTheme {
					continue
				}
				target := common.ResolveRelationshipTarget(ownerPart, rel.Target)
				inv.Bindings = append(inv.Bindings, ThemeBinding{
					OwnerType:      ownerType,
					OwnerPart:      ownerPart,
					RelationshipID: rel.ID,
					ThemePart:      target,
				})
			}
		}
		return nil
	}

	if err := collect("ppt/slideMasters/_rels/", themeOwnerSlideMaster); err != nil {
		return ThemeInventory{}, err
	}
	if err := collect("ppt/notesMasters/_rels/", themeOwnerNotesMaster); err != nil {
		return ThemeInventory{}, err
	}
	sort.Slice(inv.Bindings, func(i, j int) bool {
		if inv.Bindings[i].OwnerPart == inv.Bindings[j].OwnerPart {
			return inv.Bindings[i].RelationshipID < inv.Bindings[j].RelationshipID
		}
		return inv.Bindings[i].OwnerPart < inv.Bindings[j].OwnerPart
	})
	return inv, nil
}

func (e *PresentationEditor) ThemeInventory() (ThemeInventory, error) {
	return e.GetThemeInventory()
}

// SetThemeData replaces one concrete theme part with caller-provided XML data.
func (e *PresentationEditor) SetThemeData(partPath string, data []byte) error {
	if e == nil || e.parts == nil {
		return errors.New("editor cannot be nil")
	}
	path := strings.TrimSpace(common.CanonicalPartPath(partPath))
	if !themePartPathPattern.MatchString(path) {
		return fmt.Errorf("invalid theme part path %q", partPath)
	}
	if err := validateThemeXML(data); err != nil {
		return err
	}
	e.parts.Set(path, append([]byte(nil), data...))
	return nil
}

// SetThemeFontScheme updates the major/minor latin typefaces across all discovered themes.
func (e *PresentationEditor) SetThemeFontScheme(major, minor string) error {
	major = strings.TrimSpace(major)
	minor = strings.TrimSpace(minor)
	if major == "" || minor == "" {
		return errors.New("major and minor fonts are required")
	}
	inv, err := e.GetThemeInventory()
	if err != nil {
		return err
	}
	for _, themePath := range inv.ThemeParts {
		xmlData, ok := e.parts.Get(themePath)
		if !ok {
			return fmt.Errorf("theme part not found: %s", themePath)
		}
		rewritten, err := rewriteThemeFontScheme(xmlData, major, minor)
		if err != nil {
			return fmt.Errorf("rewrite font scheme %s: %w", themePath, err)
		}
		e.parts.Set(themePath, rewritten)
	}
	return nil
}

// SetThemeColorScheme updates the standard 12 theme color slots across all discovered themes.
func (e *PresentationEditor) SetThemeColorScheme(s ThemeColorScheme) error {
	if s.isEmpty() {
		return errors.New("at least one theme color must be provided")
	}
	inv, err := e.GetThemeInventory()
	if err != nil {
		return err
	}
	for _, themePath := range inv.ThemeParts {
		xmlData, ok := e.parts.Get(themePath)
		if !ok {
			return fmt.Errorf("theme part not found: %s", themePath)
		}
		rewritten := rewriteThemeColors(xmlData, s)
		e.parts.Set(themePath, rewritten)
	}
	return nil
}

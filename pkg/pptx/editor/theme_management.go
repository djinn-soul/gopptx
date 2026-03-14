package editor

import (
	"bytes"
	"encoding/xml"
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

func sourcePartFromRels(relsPath string) string {
	path := strings.ReplaceAll(strings.TrimSpace(relsPath), "\\", "/")
	path = strings.Replace(path, "/_rels/", "/", 1)
	return strings.TrimSuffix(path, ".rels")
}

func validateThemeXML(data []byte) error {
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.Token()
		if err != nil {
			return fmt.Errorf("invalid theme xml: %w", err)
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Local != "theme" {
			return fmt.Errorf("invalid theme xml root <%s>", start.Name.Local)
		}
		return nil
	}
}

func rewriteThemeFontScheme(xmlData []byte, major, minor string) ([]byte, error) {
	out := string(xmlData)
	var ok bool
	out, ok = rewriteFontBlock(out, majorFontPattern, major)
	if !ok {
		return nil, errors.New("missing majorFont block")
	}
	out, ok = rewriteFontBlock(out, minorFontPattern, minor)
	if !ok {
		return nil, errors.New("missing minorFont block")
	}
	return []byte(out), nil
}

func rewriteFontBlock(xmlText string, block *regexp.Regexp, typeface string) (string, bool) {
	match := block.FindString(xmlText)
	if match == "" {
		return xmlText, false
	}
	latin := latinPattern.FindString(match)
	if latin == "" {
		return xmlText, false
	}
	rewrittenLatin := latin
	if typefacePattern.MatchString(rewrittenLatin) {
		rewrittenLatin = typefacePattern.ReplaceAllString(rewrittenLatin, `typeface="`+escapeXMLAttr(typeface)+`"`)
	} else {
		rewrittenLatin = strings.Replace(rewrittenLatin, "/>", ` typeface="`+escapeXMLAttr(typeface)+`"/>`, 1)
	}
	match = strings.Replace(match, latin, rewrittenLatin, 1)
	return strings.Replace(xmlText, block.FindString(xmlText), match, 1), true
}

func rewriteThemeColors(xmlData []byte, scheme ThemeColorScheme) []byte {
	out := string(xmlData)
	for _, slot := range scheme.toSlots() {
		if slot.hex == "" {
			continue
		}
		out = rewriteColorSlot(out, slot.name, slot.hex)
	}
	return []byte(out)
}

func rewriteColorSlot(xmlText, slotName, hex string) string {
	pattern := regexp.MustCompile(`(?is)<(?:\w+:)?` + slotName + `\b[^>]*>.*?</(?:\w+:)?` + slotName + `>`)
	block := pattern.FindString(xmlText)
	if block == "" {
		return xmlText
	}
	if strings.Contains(strings.ToLower(block), "srgbclr") {
		valPattern := regexp.MustCompile(`\bval="[^"]*"`)
		if valPattern.MatchString(block) {
			block = valPattern.ReplaceAllString(block, `val="`+escapeXMLAttr(hex)+`"`)
		}
	} else {
		lastPattern := regexp.MustCompile(`\blastClr="[^"]*"`)
		if lastPattern.MatchString(block) {
			block = lastPattern.ReplaceAllString(block, `lastClr="`+escapeXMLAttr(hex)+`"`)
		}
	}
	return strings.Replace(xmlText, pattern.FindString(xmlText), block, 1)
}

func escapeXMLAttr(value string) string {
	replacer := strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;", `"`, "&quot;")
	return replacer.Replace(value)
}

type colorSlot struct {
	name string
	hex  string
}

func (s ThemeColorScheme) toSlots() []colorSlot {
	return []colorSlot{
		{name: "dk1", hex: cleanHex(s.Dk1)},
		{name: "lt1", hex: cleanHex(s.Lt1)},
		{name: "dk2", hex: cleanHex(s.Dk2)},
		{name: "lt2", hex: cleanHex(s.Lt2)},
		{name: "accent1", hex: cleanHex(s.Accent1)},
		{name: "accent2", hex: cleanHex(s.Accent2)},
		{name: "accent3", hex: cleanHex(s.Accent3)},
		{name: "accent4", hex: cleanHex(s.Accent4)},
		{name: "accent5", hex: cleanHex(s.Accent5)},
		{name: "accent6", hex: cleanHex(s.Accent6)},
		{name: "hlink", hex: cleanHex(s.Hlink)},
		{name: "folHlink", hex: cleanHex(s.FolHlink)},
	}
}

func cleanHex(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "#")
}

func (s ThemeColorScheme) isEmpty() bool {
	for _, slot := range s.toSlots() {
		if slot.hex != "" {
			return false
		}
	}
	return true
}

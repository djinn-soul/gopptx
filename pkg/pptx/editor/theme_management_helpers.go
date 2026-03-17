package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// colorSlot pairs a theme color slot name with its hex value.
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

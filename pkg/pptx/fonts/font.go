package fonts

import (
	"encoding/hex"
	"strings"
)

// FontStyle represents the style variant of an embedded font.
type FontStyle int

// fontStyleRegular is the OOXML element name for the default font style.
const fontStyleRegular = "regular"

const (
	defaultPitchFamily = 0x22
	minObfuscateBytes  = 32
	guidHexLen         = 32
	guidKeyBytes       = 16
)

const (
	// StyleRegular is the default regular font style.
	StyleRegular FontStyle = iota
	// StyleBold is the bold font style.
	StyleBold
	// StyleItalic is the italic font style.
	StyleItalic
	// StyleBoldItalic is the bold and italic font style.
	StyleBoldItalic
)

// XMLElement returns the corresponding OOXML element name for the style.
func (s FontStyle) XMLElement() string {
	switch s {
	case StyleRegular:
		return fontStyleRegular
	case StyleBold:
		return "bold"
	case StyleItalic:
		return "italic"
	case StyleBoldItalic:
		return "boldItalic"
	default:
		return fontStyleRegular
	}
}

// FontCharset represents the character set for the font.
type FontCharset uint8

const (
	CharsetAnsi        FontCharset = 0x00
	CharsetSymbol      FontCharset = 0x02
	CharsetShiftJis    FontCharset = 0x80
	CharsetHangul      FontCharset = 0x81
	CharsetGb2312      FontCharset = 0x86
	CharsetChineseBig5 FontCharset = 0x88
	CharsetGreek       FontCharset = 0xA1
	CharsetTurkish     FontCharset = 0xA2
	CharsetHebrew      FontCharset = 0xB1
	CharsetArabic      FontCharset = 0xB2
	CharsetBaltic      FontCharset = 0xBA
	CharsetRussian     FontCharset = 0xCC
	CharsetThai        FontCharset = 0xDE
	CharsetEastEurope  FontCharset = 0xEE
)

// EmbeddedFont represents a single font embedded into the presentation.
// The Data field should contain the obfuscated .fntdata bytes ready to be written.
type EmbeddedFont struct {
	Typeface       string
	Style          FontStyle
	Charset        FontCharset
	Panose         string
	PitchFamily    uint8
	Data           []byte
	RelationshipID string // Used internally during rendering
}

// New creates a new EmbeddedFont entry.
func New(typeface string, style FontStyle, data []byte) *EmbeddedFont {
	return &EmbeddedFont{
		Typeface:    typeface,
		Style:       style,
		Charset:     CharsetAnsi,
		PitchFamily: defaultPitchFamily, // Variable pitch, Roman family
		Data:        data,
	}
}

// WithCharset sets the character set for the font.
func (f *EmbeddedFont) WithCharset(charset FontCharset) *EmbeddedFont {
	f.Charset = charset
	return f
}

// WithPanose sets the Panose-1 classification number sequence.
func (f *EmbeddedFont) WithPanose(panose string) *EmbeddedFont {
	f.Panose = panose
	return f
}

// WithPitchFamily sets the pitch and family metadata for the font.
func (f *EmbeddedFont) WithPitchFamily(pitchFamily uint8) *EmbeddedFont {
	f.PitchFamily = pitchFamily
	return f
}

// ObfuscateFont implements the OpenXML font obfuscation algorithm.
// It requires the original font file data and a valid GUID string (e.g. from the presentation metadata or generated).
// The algorithm XORs the first 32 bytes of the font data with the reverse of the GUID bytes.
func ObfuscateFont(fontData []byte, guid string) []byte {
	if len(fontData) < minObfuscateBytes {
		// Font data is too short to obfuscate fully, return as-is or partial.
		// A valid TTF/OTF is much larger than 32 bytes.
		out := make([]byte, len(fontData))
		copy(out, fontData)
		return out
	}

	// Parse GUID into 16 bytes. Remove hyphens and braces if present.
	cleanGUID := strings.ReplaceAll(guid, "-", "")
	cleanGUID = strings.ReplaceAll(cleanGUID, "{", "")
	cleanGUID = strings.ReplaceAll(cleanGUID, "}", "")

	if len(cleanGUID) != guidHexLen {
		// Invalid GUID format for obfuscation, just return a copy of the origin data
		out := make([]byte, len(fontData))
		copy(out, fontData)
		return out
	}

	guidBytes, err := hex.DecodeString(cleanGUID)
	if err != nil || len(guidBytes) != guidKeyBytes {
		out := make([]byte, len(fontData))
		copy(out, fontData)
		return out
	}

	// Reverse the GUID bytes to create the 16-byte XOR key
	key := make([]byte, guidKeyBytes)
	for i := range guidKeyBytes {
		key[i] = guidBytes[15-i]
	}

	// Copy original data
	out := make([]byte, len(fontData))
	copy(out, fontData)

	// XOR the first 32 bytes with the 16-byte key (repeated twice)
	for i := range minObfuscateBytes {
		out[i] ^= key[i%len(key)]
	}

	return out
}

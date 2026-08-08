package text

import "strings"

// ScriptKind says which typeface slot of an <a:rPr> a script is rendered from.
type ScriptKind int

const (
	// ScriptLatin is rendered from <a:latin>.
	ScriptLatin ScriptKind = iota
	// ScriptEastAsian is rendered from <a:ea>.
	ScriptEastAsian
	// ScriptComplex is rendered from <a:cs>.
	ScriptComplex
)

// ISO 15924 script codes used by the language table.
const (
	ScriptCodeArabic     = "Arab"
	ScriptCodeBengali    = "Beng"
	ScriptCodeCanadian   = "Cans"
	ScriptCodeCherokee   = "Cher"
	ScriptCodeDevanagari = "Deva"
	ScriptCodeEthiopic   = "Ethi"
	ScriptCodeGeorgian   = "Geor"
	ScriptCodeGujarati   = "Gujr"
	ScriptCodeGurmukhi   = "Guru"
	ScriptCodeHangul     = "Hang"
	ScriptCodeHanSimp    = "Hans"
	ScriptCodeHanTrad    = "Hant"
	ScriptCodeHebrew     = "Hebr"
	ScriptCodeJapanese   = "Jpan"
	ScriptCodeKannada    = "Knda"
	ScriptCodeKhmer      = "Khmr"
	ScriptCodeLao        = "Laoo"
	ScriptCodeMalayalam  = "Mlym"
	ScriptCodeOriya      = "Orya"
	ScriptCodeSinhala    = "Sinh"
	ScriptCodeSyriac     = "Syrc"
	ScriptCodeTamil      = "Taml"
	ScriptCodeTelugu     = "Telu"
	ScriptCodeThaana     = "Thaa"
	ScriptCodeThai       = "Thai"
	ScriptCodeTibetan    = "Tibt"
	ScriptCodeVietnamese = "Viet"
)

// Primary language subtags that appear in more than one table below.
const (
	langArabic  = "ar"
	langPersian = "fa"
	langUrdu    = "ur"
	langPashto  = "ps"
	langKurdish = "ku"
	langSorani  = "ckb"
	langSindhi  = "sd"
	langHebrew  = "he"
	langYiddish = "yi"
	// langHebrewLegacy is the pre-1989 ISO 639 code for Hebrew, still emitted by
	// some producers.
	langHebrewLegacy = "iw"
	langDhivehi      = "dv"
	langSyriac       = "syr"
)

// langScript maps a primary language subtag to its ISO 15924 script code, using
// the mapping from the issue (docx4j's LanguageTagToScriptMapping). Chinese is
// resolved separately because it depends on the region subtag.
//
//nolint:gochecknoglobals // immutable lookup table
var langScript = map[string]string{
	"ja": ScriptCodeJapanese,
	"ko": ScriptCodeHangul,
	// The Arabic-script languages. Persian, Urdu, Pashto and Kurdish were absent
	// until they were noticed falling through to Latin, which put their text in
	// <a:latin> — the wrong typeface slot for a complex script.
	langArabic: ScriptCodeArabic, langPersian: ScriptCodeArabic, langUrdu: ScriptCodeArabic,
	langPashto: ScriptCodeArabic, langKurdish: ScriptCodeArabic, langSorani: ScriptCodeArabic,
	langHebrew: ScriptCodeHebrew, langYiddish: ScriptCodeHebrew, langHebrewLegacy: ScriptCodeHebrew,
	"th":  ScriptCodeThai,
	"ti":  ScriptCodeEthiopic,
	"bwo": ScriptCodeEthiopic, "eth": ScriptCodeEthiopic,
	"kxh": ScriptCodeEthiopic, "mdy": ScriptCodeEthiopic,
	"bn": ScriptCodeBengali, "as": ScriptCodeBengali, "mni": ScriptCodeBengali,
	"gu":        ScriptCodeGujarati,
	"km":        ScriptCodeKhmer,
	"kn":        ScriptCodeKannada,
	"pa":        ScriptCodeGurmukhi,
	"iu":        ScriptCodeCanadian,
	"chr":       ScriptCodeCherokee,
	"bo":        ScriptCodeTibetan,
	langDhivehi: ScriptCodeThaana,
	"hi":        ScriptCodeDevanagari,
	"ks":        ScriptCodeDevanagari, "kok": ScriptCodeDevanagari, "mr": ScriptCodeDevanagari,
	"ne": ScriptCodeDevanagari, "sa": ScriptCodeDevanagari, langSindhi: ScriptCodeDevanagari,
	"te":       ScriptCodeTelugu,
	"ta":       ScriptCodeTamil,
	langSyriac: ScriptCodeSyriac,
	"or":       ScriptCodeOriya,
	"ml":       ScriptCodeMalayalam,
	"lo":       ScriptCodeLao,
	"si":       ScriptCodeSinhala,
	"vi":       ScriptCodeVietnamese,
	"lha":      ScriptCodeVietnamese, "nut": ScriptCodeVietnamese,
	"ka": ScriptCodeGeorgian,
}

// Default fonts for the right-to-left languages.
const (
	rtlFontArial     = "Arial"
	rtlFontTahoma    = "Tahoma"
	rtlFontDavid     = "David"
	rtlFontNastaleeq = "Jameel Noori Nastaleeq"
	rtlFontThaana    = "MV Boli"
	rtlFontSyriac    = "Estrangelo Edessa"
)

// rtlLanguages are the languages written right to left, with the font
// PowerPoint reaches for by default. The direction is a property of the
// language, so a caller need not set ParagraphStyle.RTL by hand.
//
//nolint:gochecknoglobals // immutable lookup table
var rtlLanguages = map[string]string{
	langArabic:       rtlFontArial,
	langPersian:      rtlFontTahoma,
	langUrdu:         rtlFontNastaleeq,
	langPashto:       rtlFontTahoma,
	langKurdish:      rtlFontTahoma,
	langSorani:       rtlFontTahoma,
	langSindhi:       rtlFontArial,
	langHebrew:       rtlFontDavid,
	langYiddish:      rtlFontDavid,
	langHebrewLegacy: rtlFontDavid,
	langDhivehi:      rtlFontThaana,
	langSyriac:       rtlFontSyriac,
}

// eastAsianScripts are the scripts PowerPoint renders from the <a:ea> typeface.
// Everything else non-Latin in the table is a complex script, drawn from <a:cs>.
//
//nolint:gochecknoglobals // immutable lookup table
var eastAsianScripts = map[string]bool{
	ScriptCodeJapanese: true,
	ScriptCodeHangul:   true,
	ScriptCodeHanSimp:  true,
	ScriptCodeHanTrad:  true,
}

// LanguageToScript returns the ISO 15924 script code for a BCP-47 language tag,
// or "" when the tag is Latin-script or unrecognised (upstream issue #172).
//
// Both "-" and "_" are accepted as the subtag separator, since the mapping in
// the issue was written against tags like "zh_CN".
func LanguageToScript(langTag string) string {
	tag := strings.TrimSpace(langTag)
	if tag == "" {
		return ""
	}
	normalized := strings.ReplaceAll(tag, "_", "-")
	lang := strings.ToLower(normalized)
	region := ""
	if dash := strings.Index(lang, "-"); dash >= 0 {
		region = strings.ToLower(normalized[dash+1:])
		lang = lang[:dash]
	}

	// A tag that states its script outranks the primary language's default.
	// Kurdish is the case that exposed this: ku defaults to Arabic, so ku-Latn
	// — an unambiguous statement that this text is Latin — was answered Arab.
	if script := explicitScriptSubtag(normalized); script != "" {
		return script
	}

	if lang == "zh" {
		// Mainland China and Singapore use simplified characters; everywhere
		// else that writes Chinese uses traditional.
		if region == "cn" || region == "sg" || strings.HasPrefix(region, "hans") {
			return ScriptCodeHanSimp
		}
		return ScriptCodeHanTrad
	}
	return langScript[lang]
}

// IsRTLLanguage reports whether a BCP-47 language tag is written right to left.
func IsRTLLanguage(langTag string) bool {
	// An explicit script settles the direction on its own: ku-Latn is Kurdish
	// written in Latin letters and reads left to right, whatever the primary
	// subtag's usual script would be.
	if script := explicitScriptSubtag(strings.ReplaceAll(strings.TrimSpace(langTag), "_", "-")); script != "" {
		return rtlScripts[script]
	}
	_, ok := rtlLanguages[primaryLanguageSubtag(langTag)]
	return ok
}

// rtlScripts are the ISO 15924 codes written right to left, for tags that name
// their script outright.
//
//nolint:gochecknoglobals // immutable lookup table
var rtlScripts = map[string]bool{
	ScriptCodeArabic: true,
	ScriptCodeHebrew: true,
	ScriptCodeThaana: true,
	ScriptCodeSyriac: true,
}

// explicitScriptSubtag returns the ISO 15924 code a tag states for itself, or
// an empty string when it states none.
//
// BCP-47 puts the script second and spells it as four letters — "ku-Latn",
// "zh-Hans", "sr-Cyrl" — which is exactly the code this package uses, so it
// only has to be found and cased canonically. A two-letter second subtag is a
// region, not a script, and is left alone.
func explicitScriptSubtag(normalizedTag string) string {
	parts := strings.Split(normalizedTag, "-")
	if len(parts) < 2 {
		return ""
	}
	candidate := parts[1]
	if len(candidate) != scriptSubtagLength {
		return ""
	}
	for _, r := range candidate {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return ""
		}
	}
	return strings.ToUpper(candidate[:1]) + strings.ToLower(candidate[1:])
}

// scriptSubtagLength is the length BCP-47 gives a script subtag.
const scriptSubtagLength = 4

// RTLDefaultFont returns the font PowerPoint defaults to for a right-to-left
// language, or "" when the tag is not one.
func RTLDefaultFont(langTag string) string {
	return rtlLanguages[primaryLanguageSubtag(langTag)]
}

// primaryLanguageSubtag lower-cases a tag and drops everything after the first
// separator, accepting both "-" and "_".
func primaryLanguageSubtag(langTag string) string {
	lang := strings.ToLower(strings.TrimSpace(langTag))
	lang = strings.ReplaceAll(lang, "_", "-")
	if dash := strings.Index(lang, "-"); dash >= 0 {
		lang = lang[:dash]
	}
	return lang
}

// ScriptKindForLanguage reports which typeface slot a language tag's text is
// drawn from, so a caller can put the font in <a:latin>, <a:ea> or <a:cs>.
func ScriptKindForLanguage(langTag string) ScriptKind {
	script := LanguageToScript(langTag)
	switch {
	case script == "":
		return ScriptLatin
	case eastAsianScripts[script]:
		return ScriptEastAsian
	default:
		return ScriptComplex
	}
}

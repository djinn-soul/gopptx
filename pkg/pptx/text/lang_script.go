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

// langScript maps a primary language subtag to its ISO 15924 script code, using
// the mapping from the issue (docx4j's LanguageTagToScriptMapping). Chinese is
// resolved separately because it depends on the region subtag.
//
//nolint:gochecknoglobals // immutable lookup table
var langScript = map[string]string{
	"ja": "Jpan",
	"ko": "Hang",
	"ar": "Arab",
	"he": "Hebr", "yi": "Hebr", "iw": "Hebr",
	"th":  "Thai",
	"ti":  "Ethi",
	"bwo": "Ethi", "eth": "Ethi", "kxh": "Ethi", "mdy": "Ethi",
	"bn": "Beng", "as": "Beng", "mni": "Beng",
	"gu":  "Gujr",
	"km":  "Khmr",
	"kn":  "Knda",
	"pa":  "Guru",
	"iu":  "Cans",
	"chr": "Cher",
	"bo":  "Tibt",
	"dv":  "Thaa",
	"hi":  "Deva",
	"ks":  "Deva", "kok": "Deva", "mr": "Deva", "ne": "Deva", "sa": "Deva", "sd": "Deva",
	"te":  "Telu",
	"ta":  "Taml",
	"syr": "Syrc",
	"or":  "Orya",
	"ml":  "Mlym",
	"lo":  "Laoo",
	"si":  "Sinh",
	"vi":  "Viet",
	"lha": "Viet", "nut": "Viet",
	"ka": "Geor",
}

// eastAsianScripts are the scripts PowerPoint renders from the <a:ea> typeface.
// Everything else non-Latin in the table is a complex script, drawn from <a:cs>.
//
//nolint:gochecknoglobals // immutable lookup table
var eastAsianScripts = map[string]bool{
	"Jpan": true,
	"Hang": true,
	"Hans": true,
	"Hant": true,
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

	if lang == "zh" {
		// Mainland China and Singapore use simplified characters; everywhere
		// else that writes Chinese uses traditional.
		if region == "cn" || region == "sg" || strings.HasPrefix(region, "hans") {
			return "Hans"
		}
		return "Hant"
	}
	return langScript[lang]
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

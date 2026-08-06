package export

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/signintech/gopdf"
)

const codeTokenHintThreshold = 2

func setPDFFontAliases(pdf *gopdf.GoPdf, sansAlias, serifAlias, monoAlias string) {
	documentFontsFor(pdf).setGenericAliases(sansAlias, serifAlias, monoAlias)
}

func setPDFCJKAlias(pdf *gopdf.GoPdf, alias string) {
	documentFontsFor(pdf).setCJKAlias(alias)
}

func fallbackAlias(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func setPDFTextFontWithHint(pdf *gopdf.GoPdf, size int, bold bool, italic bool, fontHint string) {
	setPDFTextFontWithHintAndLang(pdf, size, bold, italic, fontHint, "")
}

func setPDFTextFontWithHintAndLang(
	pdf *gopdf.GoPdf,
	size int,
	bold bool,
	italic bool,
	fontHint string,
	lang string,
) {
	style := ""
	if bold {
		style += "B"
	}
	if italic {
		style += "I"
	}
	if size <= 0 {
		size = defaultFontSize
	}
	// The typeface the run names is embedded on demand, so a deck set in
	// Georgia is drawn in Georgia rather than in whichever serif the host
	// happens to fall back to.
	if alias := ensureNamedPDFFont(pdf, fontHint); alias != "" {
		_ = pdf.SetFont(alias, style, size)
		return
	}
	_ = pdf.SetFont(resolvePDFFontAliasForRun(pdf, fontHint, lang), style, size)
}

func resolvePDFFontAlias(pdf *gopdf.GoPdf, fontHint string) string {
	if alias := cachedNamedPDFFontAlias(pdf, fontHint); alias != "" {
		return alias
	}
	sans, serif, mono, _ := documentFontsFor(pdf).genericAliases()
	hint := strings.ToLower(strings.TrimSpace(fontHint))
	switch {
	case isMonospaceFontHint(hint):
		return mono
	case isSerifFontHint(hint):
		return serif
	default:
		return sans
	}
}

func resolvePDFFontAliasForRun(pdf *gopdf.GoPdf, fontHint, lang string) string {
	if hint := strings.TrimSpace(fontHint); hint != "" {
		return resolvePDFFontAlias(pdf, hint)
	}
	sans, _, _, cjk := documentFontsFor(pdf).genericAliases()
	normalizedLang := strings.ToLower(strings.TrimSpace(lang))
	switch {
	case strings.HasPrefix(normalizedLang, "ja"),
		strings.HasPrefix(normalizedLang, "zh"),
		strings.HasPrefix(normalizedLang, "ko"):
		return fallbackAlias(cjk, sans)
	default:
		return sans
	}
}

func isMonospaceFontHint(hint string) bool {
	return strings.Contains(hint, "mono") ||
		strings.Contains(hint, "consolas") ||
		strings.Contains(hint, "courier") ||
		strings.Contains(hint, "code")
}

func isSerifFontHint(hint string) bool {
	return strings.Contains(hint, "serif") ||
		strings.Contains(hint, "times") ||
		strings.Contains(hint, "cambria") ||
		strings.Contains(hint, "georgia")
}

func inferCodeFontHint(textValue string) string {
	if strings.TrimSpace(textValue) == "" {
		return ""
	}
	tokenHits := 0
	for _, token := range []string{"{", "}", "=>", "::", "func ", "return", "if ", "for ", "[]", "()"} {
		if strings.Contains(textValue, token) {
			tokenHits++
		}
	}
	if tokenHits < codeTokenHintThreshold {
		return ""
	}
	punct := 0
	total := 0
	for _, r := range textValue {
		if unicode.IsSpace(r) {
			continue
		}
		total++
		if strings.ContainsRune("{}[]();:=<>\"'`", r) {
			punct++
		}
	}
	if total == 0 {
		return ""
	}
	if float64(punct)/float64(total) >= 0.14 || utf8.RuneCountInString(textValue) >= 80 {
		return "Consolas"
	}
	return ""
}

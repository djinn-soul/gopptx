package export

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/signintech/gopdf"
)

//nolint:gochecknoglobals // Alias values are runtime-configurable for renderer-level font fallback tuning.
var (
	pdfSansAlias  = fontFamilySans
	pdfSerifAlias = fontFamilySans
	pdfMonoAlias  = fontFamilySans
	pdfCJKAlias   = fontFamilySans
)

const codeTokenHintThreshold = 2

func setPDFFontAliases(sansAlias, serifAlias, monoAlias string) {
	pdfSansAlias = fallbackAlias(sansAlias, fontFamilySans)
	pdfSerifAlias = fallbackAlias(serifAlias, pdfSansAlias)
	pdfMonoAlias = fallbackAlias(monoAlias, pdfSansAlias)
	pdfCJKAlias = fallbackAlias(pdfCJKAlias, pdfSansAlias)
}

func setPDFCJKAlias(alias string) {
	pdfCJKAlias = fallbackAlias(alias, pdfSansAlias)
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
	_ = pdf.SetFont(resolvePDFFontAliasForRun(fontHint, lang), style, size)
}

func resolvePDFFontAlias(fontHint string) string {
	hint := strings.ToLower(strings.TrimSpace(fontHint))
	switch {
	case isMonospaceFontHint(hint):
		return pdfMonoAlias
	case isSerifFontHint(hint):
		return pdfSerifAlias
	default:
		return pdfSansAlias
	}
}

func resolvePDFFontAliasForRun(fontHint, lang string) string {
	if hint := strings.TrimSpace(fontHint); hint != "" {
		return resolvePDFFontAlias(hint)
	}
	normalizedLang := strings.ToLower(strings.TrimSpace(lang))
	switch {
	case strings.HasPrefix(normalizedLang, "ja"),
		strings.HasPrefix(normalizedLang, "zh"),
		strings.HasPrefix(normalizedLang, "ko"):
		return fallbackAlias(pdfCJKAlias, pdfSansAlias)
	default:
		return pdfSansAlias
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

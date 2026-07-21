package export

import (
	"errors"

	"github.com/signintech/gopdf"
)

func configureNativePDFFont(pdf *gopdf.GoPdf, opts PDFOptions) error {
	sansAlias := ""
	if tryNativePDFFonts(pdf, opts.NativeFontPaths, fontFamilySans) {
		sansAlias = fontFamilySans
	} else if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilySans), fontFamilySans) {
		sansAlias = fontFamilySans
	}
	if sansAlias == "" {
		return errors.New("no system TTF font found; install Arial or DejaVu Sans, or specify NativeFontPaths")
	}
	// Register bold, italic, and bold+italic variants so SetFont("sans","B",...) works.
	tryNativePDFFontsWithStyle(pdf, systemFontPathsForFamily(fontFamilySansBold), fontFamilySans, gopdf.Bold)
	tryNativePDFFontsWithStyle(pdf, systemFontPathsForFamily(fontFamilySansItalic), fontFamilySans, gopdf.Italic)
	tryNativePDFFontsWithStyle(
		pdf, systemFontPathsForFamily(fontFamilySansBoldItalic), fontFamilySans, gopdf.Bold|gopdf.Italic,
	)

	serifAlias := ""
	if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilySerif), fontFamilySerif) {
		serifAlias = fontFamilySerif
		tryNativePDFFontsWithStyle(pdf, systemFontPathsForFamily(fontFamilySerifBold), fontFamilySerif, gopdf.Bold)
		tryNativePDFFontsWithStyle(pdf, systemFontPathsForFamily(fontFamilySerifItalic), fontFamilySerif, gopdf.Italic)
		tryNativePDFFontsWithStyle(
			pdf, systemFontPathsForFamily(fontFamilySerifBoldItalic), fontFamilySerif, gopdf.Bold|gopdf.Italic,
		)
	}
	monoAlias := ""
	if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilyMono), fontFamilyMono) {
		monoAlias = fontFamilyMono
		tryNativePDFFontsWithStyle(pdf, systemFontPathsForFamily(fontFamilyMonoBold), fontFamilyMono, gopdf.Bold)
		tryNativePDFFontsWithStyle(pdf, systemFontPathsForFamily(fontFamilyMonoItalic), fontFamilyMono, gopdf.Italic)
		tryNativePDFFontsWithStyle(
			pdf, systemFontPathsForFamily(fontFamilyMonoBoldItalic), fontFamilyMono, gopdf.Bold|gopdf.Italic,
		)
	}
	cjkAlias := ""
	if tryNativePDFFonts(pdf, systemFontPathsForFamily(fontFamilyCJK), fontFamilyCJK) {
		cjkAlias = fontFamilyCJK
	}
	setPDFFontAliases(sansAlias, serifAlias, monoAlias)
	setPDFCJKAlias(cjkAlias)
	return nil
}

func tryNativePDFFonts(pdf *gopdf.GoPdf, fontPaths []string, alias string) bool {
	for _, path := range fontPaths {
		if err := pdf.AddTTFFont(alias, path); err != nil {
			continue
		}
		if err := pdf.SetFont(alias, "", defaultFontSize); err == nil {
			return true
		}
	}
	return false
}

func tryNativePDFFontsWithStyle(pdf *gopdf.GoPdf, fontPaths []string, alias string, style int) {
	for _, path := range fontPaths {
		if err := pdf.AddTTFFontWithOption(alias, path, gopdf.TtfOption{Style: style}); err != nil {
			continue
		}
		// Verify the font can be set at this style.
		if err := pdf.SetFontWithStyle(alias, style, defaultFontSize); err == nil {
			return
		}
	}
}

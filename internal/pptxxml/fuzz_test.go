package pptxxml

import (
	"testing"
)

func FuzzEscape(f *testing.F) {
	f.Add("normal string")
	f.Add("<xml>&amp;</xml>")
	f.Add("a & b < c > d \" e ' f")
	f.Fuzz(func(_ *testing.T, s string) {
		_ = Escape(s)
	})
}

func FuzzContentTypes(f *testing.F) {
	f.Add(1, 1, 1, 1, 1, 1, 1, 1, 1)
	f.Fuzz(func(
		_ *testing.T,
		slideCount, chartCount, smartArtCount, customXMLCount int,
		masterCount, notesThemeIndex int,
		_ /* embeddedFontsCount */ int,
		imageExtCount int,
		_ /* commentSlidesCount */ int,
	) {
		// Limit counts to avoid OOM during fuzzing
		if slideCount < 0 || slideCount > 100 {
			slideCount = 10
		}
		if chartCount < 0 || chartCount > 100 {
			chartCount = 10
		}
		if smartArtCount < 0 || smartArtCount > 100 {
			smartArtCount = 10
		}
		if customXMLCount < 0 || customXMLCount > 100 {
			customXMLCount = 10
		}
		if masterCount < 0 || masterCount > 100 {
			masterCount = 10
		}
		if notesThemeIndex < 0 || notesThemeIndex > 100 {
			notesThemeIndex = 10
		}
		if imageExtCount < 0 || imageExtCount > 10 {
			imageExtCount = 2
		}

		exts := []string{"png", "jpg"}
		if imageExtCount > 0 {
			exts = make([]string, imageExtCount)
			for i := range imageExtCount {
				exts[i] = "png" // Use valid extensions to avoid panic
			}
		}

		notesSlides := []int{1, 2}
		commentSlides := []int{1, 2}

		_ = ContentTypes(
			slideCount,
			exts,
			chartCount,
			smartArtCount,
			notesSlides,
			true,
			customXMLCount,
			masterCount,
			notesThemeIndex,
			true,
			commentSlides,
			true,
			true,
			true,
			true,
			true,
		)
	})
}

func FuzzSectionListXML(f *testing.F) {
	f.Add("Section 1", "GUID-1", int64(1))
	f.Fuzz(func(_ *testing.T, name, guid string, slideID int64) {
		sections := []Section{
			{
				Name:     name,
				GUID:     guid,
				SlideIDs: []int64{slideID},
			},
		}
		_ = SectionListXML(sections)
	})
}

func FuzzPresentation(f *testing.F) {
	f.Add("Title", 10, int64(9144000), int64(6858000), 1, "hash", "salt", 100000)
	f.Fuzz(func(
		_ *testing.T,
		title string,
		slideCount int,
		width, height int64,
		masterCount int,
		hashData, saltData string,
		spinCount int,
	) {
		if slideCount < 0 || slideCount > 100 {
			slideCount = 10
		}
		if masterCount < 0 || masterCount > 100 {
			masterCount = 1
		}

		protection := &ProtectionInfo{
			HashAlgSID: 14,
			HashData:   hashData,
			SaltData:   saltData,
			SpinCount:  spinCount,
		}

		_ = Presentation(
			title,
			slideCount,
			true,
			width,
			height,
			masterCount,
			protection,
			nil,
			false,
			nil,
			nil,
			nil,
		)
	})
}

func FuzzEmbeddedFontsXML(f *testing.F) {
	f.Add("Arial", "regular", uint8(1), "panose", uint8(2), "rId1")
	f.Fuzz(func(
		_ *testing.T,
		typeface, style string,
		charset uint8,
		panose string,
		pitchFamily uint8,
		relID string,
	) {
		fonts := []EmbeddedFontRef{
			{
				Typeface:    typeface,
				Style:       style,
				Charset:     charset,
				Panose:      panose,
				PitchFamily: pitchFamily,
				RelID:       relID,
			},
		}
		_ = EmbeddedFontsXML(fonts)
	})
}

func FuzzRichTextRun(f *testing.F) {
	f.Add("Hello World", true, false, "sng", "none", true, false, "FF0000", "FFFF00", "Arial", 12, false, true, false)
	f.Fuzz(func(
		_ *testing.T,
		text string,
		bold, italic bool,
		underline, strikethrough string,
		subscript, superscript bool,
		color, highlight, font string,
		sizePt int,
		code, allCaps, smallCaps bool,
	) {
		run := TextRunSpec{
			Text:          text,
			Bold:          bold,
			Italic:        italic,
			Underline:     underline,
			Strikethrough: strikethrough,
			Subscript:     subscript,
			Superscript:   superscript,
			Color:         color,
			Highlight:     highlight,
			Font:          font,
			SizePt:        sizePt,
			Code:          code,
			AllCaps:       allCaps,
			SmallCaps:     smallCaps,
		}
		contentStyle := ContentStyleSpec{
			Bold:      false,
			Italic:    false,
			Underline: false,
			Color:     "000000",
			SizePt:    14,
		}
		_ = richTextRun(run, contentStyle)
	})
}

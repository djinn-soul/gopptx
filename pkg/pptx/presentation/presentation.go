package presentation

import (
	"archive/zip"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/fonts"
	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/notes"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const (
	minMasterCountWithNativeNotesTheme = 2
	protectionSaltBytes                = 16
	protectionHashAlgSIDSHA512         = 14
	guidRandomBytes                    = 16
	guidVersionMask                    = 0x0f
	guidVersionNibble                  = 0x40
	guidVariantMask                    = 0x3f
	guidVariantNibble                  = 0x80
	maxAuthorInitialRunes              = 2
	authorColorPaletteSize             = 10
	customXMLRelationshipPairCount     = 2
)

type Metadata struct {
	common.Metadata

	Theme         *styling.Theme
	Master        *elements.SlideMaster
	Masters       []*elements.SlideMaster
	NotesMaster   *elements.NotesMaster
	HandoutMaster *handout.HandoutMaster
	Sections      []Section
	RTL           bool
	VBA           *vba.VBAProject
	EmbeddedFonts []fonts.EmbeddedFont
}

type Section struct {
	Name         string
	SlideIndices []int // 0-based indices of slides in this section
}

type SlideSize = common.SlideSize

func GetSlideSize4x3() SlideSize {
	return common.GetSlideSize4x3()
}

func GetSlideSize16x9() SlideSize {
	return common.GetSlideSize16x9()
}

func WritePresentationPackage(
	zw *zip.Writer,
	meta Metadata,
	slides []elements.SlideContent,
	slideCount int,
) error {
	pw := pptxxml.NewPackageWriter()

	mediaCatalog, mediaErr := media.BuildMediaCatalog(slides, meta.NotesMaster)
	if mediaErr != nil {
		return mediaErr
	}

	chartParts := BuildChartParts(slides)
	smartArtParts := BuildSmartArtParts(slides)
	notesParts := notes.BuildRenderedNotesParts(slides)
	effectiveMasters := getEffectiveMasters(meta)
	masterCount := len(effectiveMasters)
	notesThemeIndex := getNotesThemeIndex(len(notesParts) > 0)

	authors, commentsBySlide, commentSlideIndices := prepareComments(meta, slides)

	hasVBA := meta.VBA.IsMacroEnabled()

	if err := addBasicPropertyFiles(
		pw, meta, slides, slideCount, len(notesParts), ChartPartCount(chartParts), SmartArtPartCount(smartArtParts),
		notesParts, masterCount, notesThemeIndex, mediaCatalog.ImageExtensions(),
		authors, commentSlideIndices, hasVBA,
	); err != nil {
		return err
	}
	addLayoutFiles(pw, masterCount)
	addMasterFiles(pw, effectiveMasters, mediaCatalog)
	addThemeFiles(pw, meta.Theme, masterCount)
	addNotesMasterFiles(pw, meta, masterCount, notesThemeIndex, mediaCatalog)
	addHandoutMasterFiles(pw, meta, masterCount)

	if err := writeMediaFiles(pw, mediaCatalog); err != nil {
		return err
	}
	if err := writeChartFiles(pw, chartParts); err != nil {
		return err
	}
	if err := writeSmartArtFiles(pw, smartArtParts); err != nil {
		return err
	}
	if err := notes.WriteNotesFiles(pw, notesParts); err != nil {
		return err
	}

	chartBySlide := chartPartBySlide(chartParts)
	smartArtBySlide := smartArtPartBySlide(smartArtParts)
	notesTargets := notes.TargetBySlide(notesParts)
	if err := renderSlides(
		pw,
		meta,
		slides,
		mediaCatalog,
		chartBySlide,
		smartArtBySlide,
		notesTargets,
		masterCount,
		commentsBySlide,
	); err != nil {
		return err
	}

	if len(authors) > 0 {
		pw.AddPart("ppt/commentAuthors.xml", pptxxml.CommentAuthorsXML(authors))
	}

	if err := writeCustomXMLParts(pw, meta.CustomXML); err != nil {
		return err
	}
	if err := writeVBAParts(pw, meta.VBA); err != nil {
		return err
	}

	for i, f := range meta.EmbeddedFonts {
		path := fmt.Sprintf("ppt/fonts/font%d.fntdata", i+1)
		pw.AddBinaryPart(path, f.Data)
	}

	return pw.WriteTo(zw)
}

func WritePackageFiles(
	zw *zip.Writer,
	meta Metadata,
	slides []elements.SlideContent,
	slideCount int,
) error {
	return WritePresentationPackage(zw, meta, slides, slideCount)
}

func convertShowSettings(s common.ShowSettings) *pptxxml.ShowSettings {
	if !s.Loop && s.Mode == common.ShowModePresent && !s.DisableTimings && !s.HideAnimation {
		return nil
	}
	return &pptxxml.ShowSettings{
		Loop:           s.Loop,
		Mode:           pptxxml.ShowMode(s.Mode),
		ShowScrollbar:  s.ShowScrollbar,
		DisableTimings: s.DisableTimings,
		HideAnimation:  s.HideAnimation,
	}
}

//nolint:funlen // OPC manifest assembly is intentionally centralized to keep package shape explicit.
func addBasicPropertyFiles(
	pw *pptxxml.PackageWriter,
	meta Metadata,
	_ []elements.SlideContent,
	slideCount, notesPartCount, chartPartCount, smartArtPartCount int,
	notesParts []notes.RenderedNotesPart,
	masterCount, notesThemeIndex int,
	mediaExtensions []string,
	authors []comments.Author,
	commentSlideIndices []int,
	hasVBA bool,
) error {
	hasNotes := notesPartCount > 0

	xSections, err := convertSections(meta.Sections, slideCount)
	if err != nil {
		return err
	}
	hasSections := len(xSections) > 0
	hasCommentAuthors := len(authors) > 0

	pw.AddPart("[Content_Types].xml", pptxxml.ContentTypes(
		slideCount, mediaExtensions, chartPartCount, smartArtPartCount,
		notes.SlideNumbers(notesParts), hasNotes,
		len(meta.CustomXML), masterCount, notesThemeIndex, hasSections, commentSlideIndices,
		meta.Protection.MarkAsFinal,
		meta.Protection.SignaturesEnabled,
		hasVBA,
		meta.HandoutMaster != nil,
		len(meta.EmbeddedFonts) > 0,
	))
	pw.AddPart("_rels/.rels", pptxxml.RootRelationships(meta.Protection.MarkAsFinal, meta.Protection.SignaturesEnabled))
	pw.AddPart(
		"ppt/_rels/presentation.xml.rels",
		pptxxml.PresentationRelationships(
			slideCount,
			hasNotes,
			len(meta.CustomXML),
			masterCount,
			hasSections,
			hasCommentAuthors,
			hasVBA,
			meta.HandoutMaster != nil,
			len(meta.EmbeddedFonts),
		),
	)
	var protInfo *pptxxml.ProtectionInfo
	if meta.Protection.ModifyPassword != "" {
		// PPT uses 16 bytes of salt by default.
		salt := make([]byte, protectionSaltBytes)
		if _, err := rand.Read(salt); err != nil {
			return fmt.Errorf("generate protection salt: %w", err)
		}
		spinCount := 100000
		hash := protection.HashModifyPassword(meta.Protection.ModifyPassword, salt, spinCount)
		protInfo = &pptxxml.ProtectionInfo{
			HashAlgSID: protectionHashAlgSIDSHA512,
			HashData:   hash,
			SaltData:   base64.StdEncoding.EncodeToString(salt),
			SpinCount:  spinCount,
		}
	}

	includeNotesMaster := hasNotes

	nextRid := 1 + masterCount + 1 + slideCount
	if includeNotesMaster {
		nextRid++
	}
	nextRid += len(meta.CustomXML) * customXMLRelationshipPairCount
	if hasSections {
		nextRid++
	}
	if hasCommentAuthors {
		nextRid++
	}
	if hasVBA {
		nextRid++
	}
	if meta.HandoutMaster != nil {
		nextRid++
	}

	xmlFonts := make([]pptxxml.EmbeddedFontRef, len(meta.EmbeddedFonts))
	for i, f := range meta.EmbeddedFonts {
		xmlFonts[i] = pptxxml.EmbeddedFontRef{
			Typeface:    f.Typeface,
			Style:       f.Style.XMLElement(),
			Charset:     uint8(f.Charset),
			Panose:      f.Panose,
			PitchFamily: f.PitchFamily,
			RelID:       "rId" + strconv.Itoa(nextRid),
		}
		nextRid++
	}

	pw.AddPart(
		"ppt/presentation.xml",
		pptxxml.Presentation(
			meta.Title, slideCount, hasNotes,
			meta.SlideSize.Width, meta.SlideSize.Height, masterCount,
			protInfo, xSections, meta.RTL, xmlFonts,
			convertShowSettings(meta.ShowSettings),
		),
	)

	if hasSections {
		pw.AddPart("ppt/sectionList.xml", pptxxml.SectionListXML(xSections))
	}

	if meta.Protection.MarkAsFinal {
		pw.AddPart("docProps/custom.xml", pptxxml.CustomProperties(true))
	}

	if meta.Protection.SignaturesEnabled {
		pw.AddPart("_xmlsignatures/origin.sigs", pptxxml.SignatureOrigin())
	}

	pw.AddPart("docProps/core.xml", pptxxml.CoreProperties(pptxxml.CorePropertiesInfo{
		Title: meta.Title, Subject: meta.Subject, Creator: meta.Creator, Description: meta.Description,
	}))
	pw.AddPart(
		"docProps/app.xml",
		pptxxml.AppProperties(slideCount, notesPartCount, meta.SlideSize.Width, meta.SlideSize.Height),
	)
	return nil
}

func convertSections(sections []Section, slideCount int) ([]pptxxml.Section, error) {
	if len(sections) == 0 {
		return nil, nil
	}
	out := make([]pptxxml.Section, len(sections))
	for i, s := range sections {
		ids := make([]int64, len(s.SlideIndices))
		for j, idx := range s.SlideIndices {
			if idx < 0 || idx >= slideCount {
				return nil, fmt.Errorf("section %q references slide index %d outside [0,%d)", s.Name, idx, slideCount)
			}
			ids[j] = int64(256 + 1 + idx)
		}
		guid, err := generateGUID()
		if err != nil {
			return nil, err
		}
		out[i] = pptxxml.Section{
			Name:     s.Name,
			GUID:     guid,
			SlideIDs: ids,
		}
	}
	return out, nil
}

func generateGUID() (string, error) {
	b := make([]byte, guidRandomBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random bytes for GUID: %w", err)
	}
	b[6] = (b[6] & guidVersionMask) | guidVersionNibble
	b[8] = (b[8] & guidVariantMask) | guidVariantNibble
	return fmt.Sprintf("{%08X-%04X-%04X-%04X-%012X}", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

//nolint:gocognit // Per-slide comment author/index normalization is intentionally explicit.
func prepareComments(
	meta Metadata,
	slides []elements.SlideContent,
) ([]comments.Author, map[int][]comments.Comment, []int) {
	var authors []comments.Author
	authorMap := make(map[string]int)

	slideComments := make(map[int][]comments.Comment)
	var commentSlideIndices []int

	commentDate := meta.GeneratedDate
	if commentDate.IsZero() {
		commentDate = time.Now()
	}

	for i, s := range slides {
		if len(s.Comments) == 0 {
			continue
		}
		var slideCms []comments.Comment
		for _, c := range s.Comments {
			idx, ok := authorMap[c.AuthorName]
			if !ok {
				idx = len(authors)
				authorMap[c.AuthorName] = idx

				initials := ""
				parts := strings.FieldsSeq(c.AuthorName)
				var initialsSb290 strings.Builder
				for p := range parts {
					if r, _ := utf8.DecodeRuneInString(p); r != utf8.RuneError {
						initialsSb290.WriteRune(r)
					}
				}
				initials += initialsSb290.String()
				if len([]rune(initials)) > maxAuthorInitialRunes {
					initials = string([]rune(initials)[:maxAuthorInitialRunes])
				}
				if initials == "" {
					initials = "A"
				}

				authors = append(authors, comments.Author{
					ID:         int64(idx + 1),
					Name:       c.AuthorName,
					Initials:   initials,
					ColorIndex: idx % authorColorPaletteSize,
					LastIndex:  0,
				})
			}

			authors[idx].LastIndex++

			x := c.X
			y := c.Y
			if x == 0 && y == 0 {
				x = 100000
				y = 100000
			}

			slideCms = append(slideCms, comments.Comment{
				ID:       int64(len(slideCms) + 1),
				AuthorID: authors[idx].ID,
				Text:     c.Text,
				Date:     commentDate,
				X:        x,
				Y:        y,
				Index:    authors[idx].LastIndex,
			})
		}
		slideComments[i] = slideCms
		commentSlideIndices = append(commentSlideIndices, i+1)
	}
	return authors, slideComments, commentSlideIndices
}

func addLayoutFiles(pw *pptxxml.PackageWriter, masterCount int) {
	layoutXMLs := []string{
		pptxxml.SlideLayoutTitleAndContent(),
		pptxxml.SlideLayoutTitleOnly(),
		pptxxml.SlideLayoutBlank(),
		pptxxml.SlideLayoutCenteredTitle(),
		pptxxml.SlideLayoutTitleAndBigContent(),
		pptxxml.SlideLayoutTwoColumn(),
	}
	for masterNum := 1; masterNum <= masterCount; masterNum++ {
		for i, xml := range layoutXMLs {
			idx := (masterNum-1)*len(layoutXMLs) + (i + 1)
			name := fmt.Sprintf("slideLayout%d.xml", idx)
			pw.AddPart(fmt.Sprintf("ppt/slideLayouts/%s", name), xml)
			pw.AddPart(fmt.Sprintf("ppt/slideLayouts/_rels/%s.rels", name), pptxxml.SlideLayoutRelationships(masterNum))
		}
	}
}

func addMasterFiles(pw *pptxxml.PackageWriter, masters []*elements.SlideMaster, mc *media.Catalog) {
	for i, master := range masters {
		masterNum := i + 1
		targets, refs := buildMasterImageInfo(master, mc)
		spec := mapMasterToSpec(master, refs)
		if spec != nil {
			spec.MasterIndex = masterNum
		}
		pw.AddPart(fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", masterNum), pptxxml.SlideMaster(spec))
		pw.AddPart(
			fmt.Sprintf("ppt/slideMasters/_rels/slideMaster%d.xml.rels", masterNum),
			pptxxml.SlideMasterRelationships(targets, masterNum, masterNum),
		)
	}
}

func addThemeFiles(pw *pptxxml.PackageWriter, theme *styling.Theme, masterCount int) {
	themeXML := pptxxml.Theme(mapThemeToSpec(theme))
	for i := 1; i <= masterCount; i++ {
		pw.AddPart(fmt.Sprintf("ppt/theme/theme%d.xml", i), themeXML)
	}
}

func addNotesMasterFiles(
	pw *pptxxml.PackageWriter,
	meta Metadata,
	masterCount, notesThemeIndex int,
	mc *media.Catalog,
) {
	if notesThemeIndex == 0 {
		return
	}

	var backgroundRID string
	var mediaName []string

	if meta.NotesMaster != nil && meta.NotesMaster.Background != nil &&
		meta.NotesMaster.Background.Type == elements.SlideBackgroundPicture &&
		meta.NotesMaster.Background.PictureFill != nil {
		if name, ok := mc.MediaNameForImage(*meta.NotesMaster.Background.PictureFill); ok {
			backgroundRID = "rId2"
			mediaName = []string{"../media/" + name}
		}
	}
	spec := elements.MapNotesMasterToSpec(meta.NotesMaster, backgroundRID)
	pw.AddPart("ppt/notesMasters/notesMaster1.xml", pptxxml.NotesMaster(spec))
	pw.AddPart(
		"ppt/notesMasters/_rels/notesMaster1.xml.rels",
		pptxxml.NotesMasterRelationships(notesThemeIndex, mediaName),
	)

	if notesThemeIndex > masterCount {
		pw.AddPart(fmt.Sprintf("ppt/theme/theme%d.xml", notesThemeIndex), pptxxml.Theme(mapThemeToSpec(meta.Theme)))
	}
}

func addHandoutMasterFiles(pw *pptxxml.PackageWriter, meta Metadata, masterCount int) {
	if meta.HandoutMaster == nil {
		return
	}
	pw.AddPart("ppt/handoutMasters/handoutMaster1.xml", meta.HandoutMaster.GenerateXML())
	pw.AddPart(
		"ppt/handoutMasters/_rels/handoutMaster1.xml.rels",
		handout.RelationshipsXML(masterCount),
	)
}

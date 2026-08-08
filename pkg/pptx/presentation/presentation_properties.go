package presentation

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/notes"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
)

//nolint:funlen // OPC manifest assembly is intentionally centralized to keep package shape explicit.
func addBasicPropertyFiles(
	pw *pptxxml.PackageWriter,
	meta Metadata,
	slides []elements.SlideContent,
	slideCount, notesPartCount, chartPartCount, smartArtPartCount int,
	notesParts []notes.RenderedNotesPart,
	masterCount, notesThemeIndex int,
	mediaExtensions []string,
	authors []comments.Author,
	commentSlideIndices []int,
	hasVBA bool,
	chartEmbeddingNumbers []int,
) error {
	hasNotes := notesPartCount > 0

	xSections, err := convertSections(meta.Sections, slideCount)
	if err != nil {
		return err
	}
	hasSections := len(xSections) > 0
	hasCommentAuthors := len(authors) > 0

	contentTypes := pptxxml.ContentTypes(
		slideCount, mediaExtensions, chartPartCount, smartArtPartCount,
		notes.SlideNumbers(notesParts), hasNotes,
		len(meta.CustomXML), masterCount, notesThemeIndex, hasSections, commentSlideIndices,
		meta.Protection.MarkAsFinal,
		meta.Protection.SignaturesEnabled,
		hasVBA,
		meta.HandoutMaster != nil,
		len(meta.EmbeddedFonts) > 0,
	)
	contentTypes = addInkContentTypes(contentTypes, slides)
	contentTypes = addMediaContentTypes(contentTypes, slides)
	// Each chart that ships a workbook needs a content type for it, or
	// PowerPoint refuses the package.
	for _, number := range chartEmbeddingNumbers {
		contentTypes = pptxxml.WithContentTypeOverride(
			contentTypes,
			ChartEmbeddingPartName(number),
			pptxxml.ChartEmbeddingContentType,
		)
	}
	presentationRels := pptxxml.PresentationRelationships(
		slideCount,
		hasNotes,
		len(meta.CustomXML),
		masterCount,
		hasSections,
		hasCommentAuthors,
		hasVBA,
		meta.HandoutMaster != nil,
		len(meta.EmbeddedFonts),
	)

	var protInfo *pptxxml.ProtectionInfo
	if meta.Protection.ModifyPassword != "" {
		salt := make([]byte, protectionSaltBytes)
		if _, err = rand.Read(salt); err != nil {
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
	// One presentation-level relationship per custom XML item. The itemProps
	// part each item also carries is related from that item's own .rels, not
	// from presentation.xml.rels, so reserving a pair here ran the counter past
	// every relationship written after it — the handout master's among them.
	nextRid += len(meta.CustomXML)
	if hasSections {
		nextRid++
	}
	if hasCommentAuthors {
		nextRid++
	}
	if hasVBA {
		nextRid++
	}
	handoutRelID := ""
	if meta.HandoutMaster != nil {
		// The handout master relationship lands here in the rels file, and
		// presentation.xml has to name that same rId in <p:handoutMasterIdLst>
		// or PowerPoint never reaches the part.
		handoutRelID = "rId" + strconv.Itoa(nextRid)
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

	// presProps, viewProps and tableStyles are parts PowerPoint writes into
	// every deck and its own package validator lists as required. presProps used
	// to appear only when print settings were set, and the other two never — yet
	// every table gopptx emits names a table style that only tableStyles.xml can
	// resolve.
	printSettingsXML := ""
	if meta.PrintSettings != nil {
		printSettingsXML = meta.PrintSettings.PrnPrXML()
	}
	customShows, err := convertCustomShows(meta.CustomShows, slideCount, masterCount)
	if err != nil {
		return err
	}
	// p:showPr belongs to CT_PresentationProperties, so the show settings are
	// written into presProps rather than presentation.xml.
	showPrXML := pptxxml.ShowPrXML(convertShowSettings(meta.ShowSettings, customShows))
	standardParts := []struct {
		partName    string
		content     string
		contentType string
		relType     string
		relTarget   string
	}{
		{
			pptxxml.PresPropsPartName, pptxxml.PresentationProps(printSettingsXML, showPrXML),
			pptxxml.PresPropsContentType,
			pptxxml.PresPropsRelationshipType, pptxxml.PresPropsRelationshipTarget,
		},
		{
			pptxxml.ViewPropsPartName, pptxxml.ViewProps(),
			pptxxml.ViewPropsContentType,
			pptxxml.ViewPropsRelationshipType, pptxxml.ViewPropsRelationshipTarget,
		},
		{
			pptxxml.TableStylesPartName, pptxxml.TableStyles(),
			pptxxml.TableStylesContentType,
			pptxxml.TableStylesRelationshipType, pptxxml.TableStylesRelationshipTarget,
		},
	}
	for _, part := range standardParts {
		pw.AddPart(part.partName, part.content)
		contentTypes = pptxxml.WithContentTypeOverride(contentTypes, part.partName, part.contentType)
		presentationRels = pptxxml.WithRelationship(
			presentationRels,
			"rId"+strconv.Itoa(nextRid),
			part.relType,
			part.relTarget,
		)
		nextRid++
	}

	pw.AddPart("[Content_Types].xml", contentTypes)
	pw.AddPart("_rels/.rels", pptxxml.RootRelationships(meta.Protection.MarkAsFinal, meta.Protection.SignaturesEnabled))
	pw.AddPart("ppt/_rels/presentation.xml.rels", presentationRels)

	pw.AddPart(
		"ppt/presentation.xml",
		pptxxml.Presentation(
			meta.Title, slideCount, hasNotes,
			meta.SlideSize.Width, meta.SlideSize.Height, masterCount,
			protInfo, xSections, meta.RTL, xmlFonts,
			customShows,
			handoutRelID,
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

	pw.AddPart("docProps/core.xml", pptxxml.CoreProperties(coreProperties(meta)))
	pw.AddPart("docProps/app.xml", pptxxml.AppProperties(pptxxml.AppPropertiesInfo{
		SlideCount:   slideCount,
		NotesCount:   notesPartCount,
		HiddenSlides: hiddenSlideCount(slides),
		Width:        meta.SlideSize.Width,
		Height:       meta.SlideSize.Height,
		Application:  meta.AppProperties.Application,
		AppVersion:   meta.AppProperties.AppVersion,
		Company:      meta.AppProperties.Company,
		Manager:      meta.AppProperties.Manager,
	}))
	return nil
}

// coreProperties merges the top-level metadata fields with the fuller
// CoreProperties struct, which the generator used to ignore entirely: the
// top-level ones win when both are set, since they are the older API.
func coreProperties(meta Metadata) pptxxml.CorePropertiesInfo {
	core := meta.CoreProperties
	pick := func(primary, secondary string) string {
		if primary != "" {
			return primary
		}
		return secondary
	}
	return pptxxml.CorePropertiesInfo{
		Title:          pick(meta.Title, core.Title),
		Subject:        pick(meta.Subject, core.Subject),
		Creator:        pick(meta.Creator, core.Creator),
		Description:    pick(meta.Description, core.Description),
		Keywords:       core.Keywords,
		Category:       core.Category,
		ContentStatus:  core.ContentStatus,
		Identifier:     core.Identifier,
		Language:       core.Language,
		Version:        core.Version,
		LastModifiedBy: core.LastModifiedBy,
		Revision:       core.Revision,
		LastPrinted:    core.LastPrinted,
		Created:        core.Created,
		Modified:       core.Modified,
	}
}

func hiddenSlideCount(slides []elements.SlideContent) int {
	count := 0
	for _, slide := range slides {
		if slide.Hidden {
			count++
		}
	}
	return count
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
		out[i] = pptxxml.Section{Name: s.Name, GUID: guid, SlideIDs: ids}
	}
	return out, nil
}

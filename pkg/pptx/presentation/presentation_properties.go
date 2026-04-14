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
		out[i] = pptxxml.Section{Name: s.Name, GUID: guid, SlideIDs: ids}
	}
	return out, nil
}

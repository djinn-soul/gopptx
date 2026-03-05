package slide

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const vbaProjectRelType = "http://schemas.microsoft.com/office/2006/relationships/vbaProject"

func RenderPresentationRelsXML(
	nonSlide []common.EditorRelationship,
	slides []common.EditorSlideRef,
	hasSections bool,
	hasVBA bool,
) (string, error) {
	rels, used, hasSectionRel, err := collectNonSlideRelationships(nonSlide, len(slides))
	if err != nil {
		return "", err
	}
	if hasSections && !hasSectionRel {
		rels = append(rels, makeSectionListRelationship(rels, slides))
	}
	if hasVBA {
		hasVBARel := false
		for _, r := range rels {
			if r.Type == vbaProjectRelType {
				hasVBARel = true
				break
			}
		}
		if !hasVBARel {
			rels = append(rels, makeVBARelationship(rels, slides))
		}
	}
	rels = appendMissingSlideRelationships(rels, used, slides)
	return RenderRelationshipsXML(rels), nil
}

func RenderRelationshipsXML(rels []common.EditorRelationship) string {
	sortRelationshipsByID(rels)
	return relationshipsXMLDocument(rels)
}

func collectNonSlideRelationships(
	nonSlide []common.EditorRelationship,
	slideCapacity int,
) ([]common.EditorRelationship, map[string]struct{}, bool, error) {
	rels := make([]common.EditorRelationship, 0, len(nonSlide)+slideCapacity+1)
	used := map[string]struct{}{}
	hasSectionRel := false
	for _, rel := range nonSlide {
		id := strings.TrimSpace(rel.ID)
		if id == "" {
			return nil, nil, false, errors.New("non-slide relationship has empty Id")
		}
		if _, exists := used[id]; exists {
			return nil, nil, false, fmt.Errorf("duplicate relationship Id %q", id)
		}
		used[id] = struct{}{}
		rels = append(rels, rel)
		if rel.Type == common.RelTypeSectionList {
			hasSectionRel = true
		}
	}
	return rels, used, hasSectionRel, nil
}

func makeSectionListRelationship(
	rels []common.EditorRelationship,
	slides []common.EditorSlideRef,
) common.EditorRelationship {
	maxNum := maxRelationshipNumber(rels, slides)
	return common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", maxNum+1),
		Type:   common.RelTypeSectionList,
		Target: "sectionList.xml",
	}
}

func makeVBARelationship(
	rels []common.EditorRelationship,
	slides []common.EditorSlideRef,
) common.EditorRelationship {
	maxNum := maxRelationshipNumber(rels, slides)
	return common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", maxNum+1),
		Type:   vbaProjectRelType,
		Target: "vbaProject.bin",
	}
}

func maxRelationshipNumber(rels []common.EditorRelationship, slides []common.EditorSlideRef) int {
	maxNum := 0
	for _, r := range rels {
		if n, ok := common.ParseRelationshipNumber(r.ID); ok && n > maxNum {
			maxNum = n
		}
	}
	for _, slide := range slides {
		if n, ok := common.ParseRelationshipNumber(slide.RelID); ok && n > maxNum {
			maxNum = n
		}
	}
	return maxNum
}

func appendMissingSlideRelationships(
	rels []common.EditorRelationship,
	used map[string]struct{},
	slides []common.EditorSlideRef,
) []common.EditorRelationship {
	for _, slide := range slides {
		if _, exists := used[slide.RelID]; exists {
			continue
		}
		rels = append(rels, common.EditorRelationship{
			ID:     slide.RelID,
			Type:   common.RelTypeSlide,
			Target: slide.Target,
		})
	}
	return rels
}

func sortRelationshipsByID(rels []common.EditorRelationship) {
	sort.Slice(rels, func(i, j int) bool {
		a, aok := common.ParseRelationshipNumber(rels[i].ID)
		b, bok := common.ParseRelationshipNumber(rels[j].ID)
		if aok && bok && a != b {
			return a < b
		}
		return rels[i].ID < rels[j].ID
	})
}

func relationshipsXMLDocument(rels []common.EditorRelationship) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString("\n")
	b.WriteString(`<Relationships xmlns="` + common.RelationshipsXMLNS + `">`)
	for _, rel := range rels {
		writeRelationshipXML(&b, rel)
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}

func writeRelationshipXML(b *strings.Builder, rel common.EditorRelationship) {
	b.WriteString("\n<Relationship Id=\"")
	b.WriteString(common.XMLEscape(rel.ID))
	b.WriteString("\" Type=\"")
	b.WriteString(common.XMLEscape(rel.Type))
	b.WriteString("\" Target=\"")
	b.WriteString(common.XMLEscape(rel.Target))
	b.WriteString("\"")
	if strings.TrimSpace(rel.TargetMode) != "" {
		b.WriteString(` TargetMode="` + common.XMLEscape(rel.TargetMode) + `"`)
	}
	b.WriteString("/>")
}

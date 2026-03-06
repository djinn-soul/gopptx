package editor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// RemoveComment removes a specific comment identified by author and author-index.
func (e *PresentationEditor) RemoveComment(slideIndex int, authorID int64, authorIndex int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	relsPath := common.SlideRelsPartName(e.slides[slideIndex].Part)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return errors.New("no relationships found for slide")
	}

	var rels relationships
	if err := xml.Unmarshal(relsData, &rels); err != nil {
		return fmt.Errorf("parse rels: %w", err)
	}

	var commentPartPath string
	for _, r := range rels.Rels {
		if r.Type == commentsRelType {
			slideDir := path.Dir(e.slides[slideIndex].Part)
			commentPartPath = path.Join(slideDir, r.Target)
			commentPartPath = strings.ReplaceAll(commentPartPath, "\\", "/")
			break
		}
	}

	if commentPartPath == "" {
		return errors.New("no comments part found for slide")
	}

	cmData, ok := e.parts.Get(commentPartPath)
	if !ok {
		return errors.New("comments part not found in store")
	}

	var lst cmLst
	if err := xml.Unmarshal(cmData, &lst); err != nil {
		return fmt.Errorf("parse comments xml: %w", err)
	}

	newComments := make([]cm, 0, len(lst.Comments))
	found := false
	for _, c := range lst.Comments {
		if c.AuthorID == authorID && c.Idx == authorIndex {
			found = true
			continue
		}
		newComments = append(newComments, c)
	}

	if !found {
		return fmt.Errorf("comment not found (authorID=%d, idx=%d)", authorID, authorIndex)
	}

	domainComments := make([]comments.Comment, len(newComments))
	for i, c := range newComments {
		t, err := parseCommentTimestamp(c.Dt)
		if err != nil {
			return fmt.Errorf("parse existing comment timestamp %q: %w", c.Dt, err)
		}
		domainComments[i] = comments.Comment{
			AuthorID: c.AuthorID,
			Text:     c.Text,
			Date:     t,
			X:        c.Pos.X,
			Y:        c.Pos.Y,
			Index:    c.Idx,
		}
	}

	xmlContent := pptxxml.CommentsXML(domainComments)
	e.parts.Set(commentPartPath, []byte(xmlContent))
	return nil
}

func (e *PresentationEditor) nextRelID(rels []relationship) int {
	maxID := 0
	for _, r := range rels {
		if n, ok := common.ParseRelationshipNumber(r.ID); ok && n > maxID {
			maxID = n
		}
	}
	return maxID + 1
}

func (e *PresentationEditor) nextCommentPartName() string {
	keys := e.parts.Keys()
	maxID := 0
	for _, p := range keys {
		if strings.HasPrefix(p, "ppt/comments/comment") && strings.HasSuffix(p, ".xml") {
			base := strings.TrimPrefix(p, "ppt/comments/comment")
			base = strings.TrimSuffix(base, ".xml")
			n, _ := strconv.Atoi(base)
			if n > maxID {
				maxID = n
			}
		}
	}
	return fmt.Sprintf("comment%d.xml", maxID+1)
}

package editor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Internal XML structures for relationships (for read/modify).
type relationships struct {
	XMLName xml.Name       `xml:"Relationships"`
	Xmlns   string         `xml:"xmlns,attr"`
	Rels    []relationship `xml:"Relationship"`
}

type relationship struct {
	ID         string `xml:"Id,attr"`
	Type       string `xml:"Type,attr"`
	Target     string `xml:"Target,attr"`
	TargetMode string `xml:"TargetMode,attr,omitempty"`
}

const (
	relsNamespace    = "http://schemas.openxmlformats.org/package/2006/relationships"
	commentsRelType  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments"
	commentsPartType = "application/vnd.openxmlformats-officedocument.presentationml.comments+xml"
)

// Internal XML structures for comments (for unmarshaling).
type cmLst struct {
	Comments []cm `xml:"cm"`
}

type cm struct {
	AuthorID int64  `xml:"authorId,attr"`
	Dt       string `xml:"dt,attr"`
	Idx      int    `xml:"idx,attr"`
	Pos      cmPos  `xml:"pos"`
	Text     string `xml:"text"`
}

type cmPos struct {
	X int64 `xml:"x,attr"`
	Y int64 `xml:"y,attr"`
}

func commentTimestampLayouts() []string {
	return []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339Nano,
		time.RFC3339,
	}
}

func parseCommentTimestamp(value string) (time.Time, error) {
	ts := strings.TrimSpace(value)
	if ts == "" {
		return time.Time{}, errors.New("empty comment timestamp")
	}
	for _, layout := range commentTimestampLayouts() {
		if t, err := time.Parse(layout, ts); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported comment timestamp %q", value)
}

// GetComments returns all comments for a specific slide.
func (e *PresentationEditor) GetComments(slideIndex int) ([]comments.Comment, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	relsPath := common.SlideRelsPartName(e.slides[slideIndex].Part)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return nil, nil // No relationships -> no comments
	}

	var rels relationships
	if err := xml.Unmarshal(relsData, &rels); err != nil {
		return nil, fmt.Errorf("parse slide rels: %w", err)
	}

	var commentPartPath string
	for _, r := range rels.Rels {
		if r.Type == commentsRelType {
			// Resolve target relative to slide part
			slideDir := path.Dir(e.slides[slideIndex].Part)
			commentPartPath = path.Join(slideDir, r.Target)
			commentPartPath = strings.ReplaceAll(commentPartPath, "\\", "/")
			break
		}
	}

	if commentPartPath == "" {
		return nil, nil
	}

	cmData, ok := e.parts.Get(commentPartPath)
	if !ok {
		return nil, nil
	}

	var lst cmLst
	if err := xml.Unmarshal(cmData, &lst); err != nil {
		return nil, fmt.Errorf("parse comments xml: %w", err)
	}

	out := make([]comments.Comment, len(lst.Comments))
	for i, c := range lst.Comments {
		t, err := parseCommentTimestamp(c.Dt)
		if err != nil {
			return nil, fmt.Errorf("parse comment timestamp %q: %w", c.Dt, err)
		}

		out[i] = comments.Comment{
			AuthorID: c.AuthorID,
			Text:     c.Text,
			Date:     t,
			X:        c.Pos.X,
			Y:        c.Pos.Y,
			Index:    c.Idx,
		}
	}
	return out, nil
}

// AddComment adds a new comment to the specified slide.
// AddComment adds a new comment to the specified slide.
func (e *PresentationEditor) AddComment(slideIndex int, authorID int64, text string, x, y int64) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	if err := e.ensureAuthorExists(authorID); err != nil {
		return err
	}

	commentPartPath, err := e.ensureCommentPart(slideIndex)
	if err != nil {
		return err
	}

	return e.saveSlideComment(commentPartPath, authorID, text, x, y)
}

func (e *PresentationEditor) ensureAuthorExists(authorID int64) error {
	e.authorCacheMu.RLock()
	_, ok := e.authorCache[authorID]
	e.authorCacheMu.RUnlock()
	if ok {
		return nil
	}

	// Try loading
	_, _ = e.GetAuthors()
	e.authorCacheMu.RLock()
	_, ok = e.authorCache[authorID]
	e.authorCacheMu.RUnlock()
	if !ok {
		return fmt.Errorf("author ID %d not found", authorID)
	}
	return nil
}

func (e *PresentationEditor) ensureCommentPart(slideIndex int) (string, error) {
	relsPath := common.SlideRelsPartName(e.slides[slideIndex].Part)
	relsData, ok := e.parts.Get(relsPath)

	var rels relationships
	if ok {
		if err := xml.Unmarshal(relsData, &rels); err != nil {
			return "", fmt.Errorf("parse rels: %w", err)
		}
	} else {
		rels.Xmlns = relsNamespace
	}

	for _, r := range rels.Rels {
		if r.Type == commentsRelType {
			slideDir := path.Dir(e.slides[slideIndex].Part)
			commentPartPath := path.Join(slideDir, r.Target)
			return strings.ReplaceAll(commentPartPath, "\\", "/"), nil
		}
	}

	newPartName := e.nextCommentPartName()
	commentPartPath := "ppt/comments/" + newPartName

	nextRelID := e.nextRelID(rels.Rels)
	rels.Rels = append(rels.Rels, relationship{
		ID:     fmt.Sprintf("rId%d", nextRelID),
		Type:   commentsRelType,
		Target: "../comments/" + newPartName,
	})

	editorRels := make([]common.EditorRelationship, len(rels.Rels))
	for i, r := range rels.Rels {
		editorRels[i] = common.EditorRelationship{
			ID: r.ID, Type: r.Type, Target: r.Target, TargetMode: r.TargetMode,
		}
	}
	e.parts.Set(relsPath, []byte(renderRelationshipsXML(editorRels)))
	return commentPartPath, nil
}

func (e *PresentationEditor) saveSlideComment(commentPartPath string, authorID int64, text string, x, y int64) error {
	var lst cmLst
	if data, ok := e.parts.Get(commentPartPath); ok {
		_ = xml.Unmarshal(data, &lst)
	}

	author := e.authorCache[authorID]
	newIdx := author.LastIndex + 1
	e.updateAuthorLastIndex(authorID, newIdx)

	now := time.Now().Format("2006-01-02T15:04:05.000")
	lst.Comments = append(lst.Comments, cm{
		AuthorID: authorID,
		Dt:       now,
		Idx:      newIdx,
		Pos:      cmPos{X: x, Y: y},
		Text:     text,
	})

	domainComments := make([]comments.Comment, len(lst.Comments))
	for i, c := range lst.Comments {
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

	e.parts.Set(commentPartPath, []byte(pptxxml.CommentsXML(domainComments)))
	return nil
}

// RemoveComment removes a comment by ID/Index?
// The task says "Remove comment by index or ID".
// Comments don't have global IDs, they have (AuthorID, Idx) tuple which is unique?
// Or we can just remove by index in the list for simplicity if that's what user expects methods to be.
// But array index is unstable.
// Let's support removing by AuthorID + Index tuple.

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

	// Filter
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

	// Save back
	// We convert to domain objects to reuse the generator
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
		if n, ok := common.ParseRelationshipNumber(r.ID); ok {
			if n > maxID {
				maxID = n
			}
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

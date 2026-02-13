package editor

import (
	"encoding/xml"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
)

const commentAuthorsPartName = "ppt/commentAuthors.xml"

// cmAuthorLst is the root element container for authors.
type cmAuthorLst struct {
	Authors []cmAuthor `xml:"cmAuthor"`
}

// cmAuthor maps to the <p:cmAuthor> element.
type cmAuthor struct {
	ID       int64  `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Initials string `xml:"initials,attr"`
	LastIdx  int    `xml:"lastIdx,attr"`
	ClrIdx   int    `xml:"clrIdx,attr"`
}

// ensureAuthorsLoadedLocked reads the commentAuthors.xml part if cache is empty.
func (e *PresentationEditor) ensureAuthorsLoadedLocked() error {
	if e.authorCache != nil {
		return nil
	}
	parsedCache := make(map[int64]comments.Author)
	nextAuthorID := int64(1)
	content, ok := e.parts.Get(commentAuthorsPartName)
	if !ok {
		// No authors file yet, that's fine.
		e.authorCache = parsedCache
		e.nextAuthorID = 1
		return nil
	}

	var lst cmAuthorLst
	if err := xml.Unmarshal(content, &lst); err != nil {
		// If unmarshal fails, we might still proceed with empty list or error out.
		// Let's error out to be safe.
		return fmt.Errorf("parse commentAuthors.xml: %w", err)
	}

	for _, ca := range lst.Authors {
		parsedCache[ca.ID] = comments.Author{
			ID:         ca.ID,
			Name:       ca.Name,
			Initials:   ca.Initials,
			LastIndex:  ca.LastIdx,
			ColorIndex: ca.ClrIdx,
		}
		if ca.ID >= nextAuthorID {
			nextAuthorID = ca.ID + 1
		}
	}
	e.authorCache = parsedCache
	e.nextAuthorID = nextAuthorID
	return nil
}

// GetAuthors returns a snapshot of all registered authors.
func (e *PresentationEditor) GetAuthors() ([]comments.Author, error) {
	e.authorCacheMu.Lock()
	defer e.authorCacheMu.Unlock()

	if err := e.ensureAuthorsLoadedLocked(); err != nil {
		return nil, err
	}

	out := make([]comments.Author, 0, len(e.authorCache))
	for _, a := range e.authorCache {
		out = append(out, a)
	}
	return out, nil
}

// AddAuthor registers a new author or returns an existing one if name+initials match.
// If an exact match is found, it is returned. Otherwise a new ID is allocated.
func (e *PresentationEditor) AddAuthor(name, initials string) (comments.Author, error) {
	e.authorCacheMu.Lock()
	defer e.authorCacheMu.Unlock()

	if err := e.ensureAuthorsLoadedLocked(); err != nil {
		return comments.Author{}, err
	}

	// Check for existing author
	for _, a := range e.authorCache {
		if a.Name == name && a.Initials == initials {
			return a, nil
		}
	}

	// Create new
	id := e.nextAuthorID
	e.nextAuthorID++

	// Simplistic color cycle: 0-based index.
	// PowerPoint usually cycles through colors.
	colorIdx := int(id) % 7 // Arbitrary cycle

	newAuthor := comments.Author{
		ID:         id,
		Name:       name,
		Initials:   initials,
		LastIndex:  0,
		ColorIndex: colorIdx,
	}
	e.authorCache[id] = newAuthor
	return newAuthor, nil
}

// updateAuthorLastIndex updates the tracking index for an author.
// This is used when adding comments to ensure unique indices.
func (e *PresentationEditor) updateAuthorLastIndex(authorID int64, newIndex int) {
	e.authorCacheMu.Lock()
	defer e.authorCacheMu.Unlock()

	if a, ok := e.authorCache[authorID]; ok {
		if newIndex > a.LastIndex {
			a.LastIndex = newIndex
			e.authorCache[authorID] = a
		}
	}
}

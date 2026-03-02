package pptx

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

// Presentation provides a high-level API for opening, modifying, and saving
// existing PPTX presentations. It wraps a PresentationEditor and simplifies
// common operations like metadata access and saving.
//
// This API is designed to be similar to python-pptx's Presentation(pptx_path)
// constructor, providing a straightforward interface for working with existing
// presentations.
type Presentation struct {
	editor *editor.PresentationEditor
	path   string // Original path used when opening
}

// Open opens an existing PPTX file for reading and modification.
// This is the primary entry point for working with existing presentations.
//
// The path parameter should be the path to an existing .pptx file.
// Returns an error if the file cannot be opened or is not a valid PPTX.
//
// Example:
//
//	prs, err := pptx.Open("existing.pptx")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer prs.Close()
func Open(path string) (*Presentation, error) {
	ed, err := editor.OpenPresentationEditor(path)
	if err != nil {
		return nil, err
	}
	return &Presentation{
		editor: ed,
		path:   path,
	}, nil
}

// Close releases any resources held by the presentation.
// It is important to call Close when done with a Presentation to avoid
// leaking file handles or other resources.
//
// After Close is called, the Presentation should not be used further.
func (p *Presentation) Close() error {
	if p == nil || p.editor == nil {
		return nil
	}
	return p.editor.Close()
}

// Save writes the presentation back to its original path.
// This method modifies the file that was opened with Open.
//
// Returns an error if the file cannot be written.
func (p *Presentation) Save() error {
	if p == nil {
		return fmt.Errorf("presentation is nil")
	}
	if p.editor == nil {
		return fmt.Errorf("presentation editor is not initialized")
	}
	return p.editor.Save(p.path)
}

// SaveAs writes the presentation to a new file path.
// This method creates a new file without modifying the original.
//
// The path parameter specifies where to save the presentation.
// Returns an error if the file cannot be written.
func (p *Presentation) SaveAs(path string) error {
	if p == nil {
		return fmt.Errorf("presentation is nil")
	}
	if p.editor == nil {
		return fmt.Errorf("presentation editor is not initialized")
	}
	return p.editor.Save(path)
}

// SlideCount returns the number of slides in the presentation.
func (p *Presentation) SlideCount() int {
	if p == nil || p.editor == nil {
		return 0
	}
	return p.editor.SlideCount()
}

// Validate performs a structural validation check on the presentation.
// Returns a list of issues found, or an empty slice if the presentation
// appears valid.
func (p *Presentation) Validate() []structural.Issue {
	if p == nil || p.editor == nil {
		return nil
	}
	return p.editor.Validate()
}

// CoreProperties returns the presentation's core properties (Dublin Core metadata).
// This provides direct access to all metadata fields in the CoreProperties struct.
func (p *Presentation) CoreProperties() common.CoreProperties {
	if p == nil || p.editor == nil {
		return common.CoreProperties{}
	}
	return p.editor.GetCoreProperties()
}

// SetCoreProperties updates all core properties at once.
func (p *Presentation) SetCoreProperties(props common.CoreProperties) {
	if p == nil || p.editor == nil {
		return
	}
	p.editor.SetCoreProperties(props)
}

// Title returns the presentation's title.
func (p *Presentation) Title() string {
	return p.CoreProperties().Title
}

// SetTitle updates the presentation's title.
func (p *Presentation) SetTitle(title string) {
	props := p.CoreProperties()
	props.Title = title
	p.SetCoreProperties(props)
}

// Subject returns the presentation's subject.
func (p *Presentation) Subject() string {
	return p.CoreProperties().Subject
}

// SetSubject updates the presentation's subject.
func (p *Presentation) SetSubject(subject string) {
	props := p.CoreProperties()
	props.Subject = subject
	p.SetCoreProperties(props)
}

// Creator returns the presentation's creator/author.
func (p *Presentation) Creator() string {
	return p.CoreProperties().Creator
}

// SetCreator updates the presentation's creator/author.
func (p *Presentation) SetCreator(creator string) {
	props := p.CoreProperties()
	props.Creator = creator
	p.SetCoreProperties(props)
}

// Author is an alias for Creator, matching python-pptx's API.
// It returns the presentation's author.
func (p *Presentation) Author() string {
	return p.Creator()
}

// SetAuthor is an alias for SetCreator, matching python-pptx's API.
// It updates the presentation's author.
func (p *Presentation) SetAuthor(author string) {
	p.SetCreator(author)
}

// Keywords returns the presentation's keywords.
func (p *Presentation) Keywords() string {
	return p.CoreProperties().Keywords
}

// SetKeywords updates the presentation's keywords.
func (p *Presentation) SetKeywords(keywords string) {
	props := p.CoreProperties()
	props.Keywords = keywords
	p.SetCoreProperties(props)
}

// Description returns the presentation's description.
func (p *Presentation) Description() string {
	return p.CoreProperties().Description
}

// SetDescription updates the presentation's description.
func (p *Presentation) SetDescription(description string) {
	props := p.CoreProperties()
	props.Description = description
	p.SetCoreProperties(props)
}

// LastModifiedBy returns the name of the person who last modified the
// presentation.
func (p *Presentation) LastModifiedBy() string {
	return p.CoreProperties().LastModifiedBy
}

// SetLastModifiedBy updates the last modified by field.
func (p *Presentation) SetLastModifiedBy(lastModifiedBy string) {
	props := p.CoreProperties()
	props.LastModifiedBy = lastModifiedBy
	p.SetCoreProperties(props)
}

// Revision returns the presentation's revision number.
func (p *Presentation) Revision() string {
	return p.CoreProperties().Revision
}

// SetRevision updates the presentation's revision number.
func (p *Presentation) SetRevision(revision string) {
	props := p.CoreProperties()
	props.Revision = revision
	p.SetCoreProperties(props)
}

// Created returns the timestamp when the presentation was created.
// The format is an ISO 8601 date-time string.
func (p *Presentation) Created() string {
	return p.CoreProperties().Created
}

// SetCreated updates the created timestamp.
// The format should be an ISO 8601 date-time string.
func (p *Presentation) SetCreated(created string) {
	props := p.CoreProperties()
	props.Created = created
	p.SetCoreProperties(props)
}

// Modified returns the timestamp when the presentation was last modified.
// The format is an ISO 8601 date-time string.
func (p *Presentation) Modified() string {
	return p.CoreProperties().Modified
}

// SetModified updates the modified timestamp.
// The format should be an ISO 8601 date-time string.
func (p *Presentation) SetModified(modified string) {
	props := p.CoreProperties()
	props.Modified = modified
	p.SetCoreProperties(props)
}

// Category returns the presentation's category.
func (p *Presentation) Category() string {
	return p.CoreProperties().Category
}

// SetCategory updates the presentation's category.
func (p *Presentation) SetCategory(category string) {
	props := p.CoreProperties()
	props.Category = category
	p.SetCoreProperties(props)
}

// ContentStatus returns the presentation's content status.
func (p *Presentation) ContentStatus() string {
	return p.CoreProperties().ContentStatus
}

// SetContentStatus updates the presentation's content status.
func (p *Presentation) SetContentStatus(status string) {
	props := p.CoreProperties()
	props.ContentStatus = status
	p.SetCoreProperties(props)
}

package pptxxml

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
)

// CommentAuthorsXML renders ppt/commentAuthors.xml.
func CommentAuthorsXML(authors []comments.Author) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:cmAuthorLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`)
	for _, author := range authors {
		b.WriteString(fmt.Sprintf(`
<p:cmAuthor id="%d" name="%s" initials="%s" lastIdx="%d" clrIdx="%d"/>`,
			author.ID, Escape(author.Name), Escape(author.Initials), author.LastIndex, author.ColorIndex))
	}
	b.WriteString(`
</p:cmAuthorLst>`)
	return b.String()
}

// CommentsXML renders ppt/comments/commentN.xml.
func CommentsXML(slideComments []comments.Comment) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:cmLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`)

	for _, cm := range slideComments {
		dt := cm.Date.Format("2006-01-02T15:04:05.000") // ISO 8601
		b.WriteString(fmt.Sprintf(`
<p:cm authorId="%d" dt="%s" idx="%d">
<p:pos x="%d" y="%d"/>
<p:text>%s</p:text>
</p:cm>`,
			cm.AuthorID, dt, cm.Index, cm.X, cm.Y, Escape(cm.Text)))
	}

	b.WriteString(`
</p:cmLst>`)
	return b.String()
}

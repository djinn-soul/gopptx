package pptxxml

import (
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
)

// CommentAuthorsXML renders ppt/commentAuthors.xml.
func CommentAuthorsXML(authors []comments.Author) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:cmAuthorLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`)
	for _, author := range authors {
		b.WriteString(`
<p:cmAuthor id="`)
		b.WriteString(strconv.FormatInt(int64(author.ID), 10))
		b.WriteString(`" name="`)
		b.WriteString(Escape(author.Name))
		b.WriteString(`" initials="`)
		b.WriteString(Escape(author.Initials))
		b.WriteString(`" lastIdx="`)
		b.WriteString(strconv.FormatInt(int64(author.LastIndex), 10))
		b.WriteString(`" clrIdx="`)
		b.WriteString(strconv.FormatInt(int64(author.ColorIndex), 10))
		b.WriteString(`"/>`)
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
		b.WriteString(`
<p:cm authorId="`)
		b.WriteString(strconv.FormatInt(int64(cm.AuthorID), 10))
		b.WriteString(`" dt="`)
		b.WriteString(dt)
		b.WriteString(`" idx="`)
		b.WriteString(strconv.FormatInt(int64(cm.Index), 10))
		b.WriteString(`">
<p:pos x="`)
		b.WriteString(strconv.FormatInt(int64(cm.X), 10))
		b.WriteString(`" y="`)
		b.WriteString(strconv.FormatInt(int64(cm.Y), 10))
		b.WriteString(`"/>
<p:text>`)
		b.WriteString(Escape(cm.Text))
		b.WriteString(`</p:text>
</p:cm>`)
	}

	b.WriteString(`
</p:cmLst>`)
	return b.String()
}

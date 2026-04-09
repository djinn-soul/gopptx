package editor

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

const textBodyRootXML = `<p:txBody xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`

// renderTextBodyXML constructs the <p:txBody> node based on Text or Runs.
// If a PresentationEditor is provided, it will register any hyperlink relationships.
func renderTextBodyXML(e *PresentationEditor, partPath string, s *parsedShape) ([]byte, error) {
	var txBody bytes.Buffer
	txBody.WriteString(textBodyRootXML)
	if err := writeTextBodyHeader(&txBody, s); err != nil {
		return nil, err
	}
	if len(s.Runs) > 0 {
		if err := writeRunParagraphXML(&txBody, e, partPath, s); err != nil {
			return nil, err
		}
	} else if err := writePlainTextParagraphXML(&txBody, s); err != nil {
		return nil, err
	}
	txBody.WriteString(`</p:txBody>`)
	return txBody.Bytes(), nil
}

func writeTextBodyHeader(txBody *bytes.Buffer, s *parsedShape) error {
	if s.TextFrame != nil {
		tfSpec, err := editorTextFrameToSpec(s.TextFrame)
		if err != nil {
			return err
		}
		txBody.WriteString(pptxxml.TextBodyPrXML(&tfSpec))
	} else {
		txBody.WriteString(`<a:bodyPr/>`)
	}
	txBody.WriteString(`<a:lstStyle/>`)
	return nil
}

func writeRunParagraphXML(
	txBody *bytes.Buffer,
	e *PresentationEditor,
	partPath string,
	s *parsedShape,
) error {
	paragraphXML, err := renderParagraphPropsXML(s.Paragraph)
	if err != nil {
		return err
	}
	txBody.WriteString(`<a:p>`)
	if paragraphXML != "" {
		txBody.WriteString(paragraphXML)
	}
	for _, r := range s.Runs {
		runSpec, err := e.editorRunToXMLSpec(partPath, r)
		if err != nil {
			return err
		}
		txBody.WriteString(normalizeEditorRunXML(pptxxml.RichTextRunXML(runSpec, pptxxml.ContentStyleSpec{})))
	}
	txBody.WriteString(`</a:p>`)
	return nil
}

func normalizeEditorRunXML(runXML string) string {
	// Preserve editor run-attribute parity when using shared run emitter.
	runXML = strings.ReplaceAll(runXML, ` cap="all"`, ` caps="all"`)
	return strings.ReplaceAll(runXML, ` cap="small"`, ` smCaps="1"`)
}

func writePlainTextParagraphXML(txBody *bytes.Buffer, s *parsedShape) error {
	paragraphXML, err := renderParagraphPropsXML(s.Paragraph)
	if err != nil {
		return err
	}
	txBody.WriteString(`<a:p>`)
	if paragraphXML != "" {
		txBody.WriteString(paragraphXML)
	}
	_, _ = fmt.Fprintf(txBody, `<a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r>`, escapeEditorText(s.Text))
	txBody.WriteString(`</a:p>`)
	return nil
}

func escapeEditorText(str string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(str))
	return buf.String()
}

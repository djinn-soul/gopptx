package common

import "strings"

// XMLElement is one depth-0 element of an XML fragment.
type XMLElement struct {
	// Name is the qualified tag name, e.g. "a:solidFill".
	Name string
	// Start is the byte offset of the element's opening '<' in the fragment.
	Start int
	// XML is the element's full text, including its children.
	XML string
}

// SplitTopLevelXMLElements lists the depth-0 elements of an XML fragment.
//
// It exists because several schema types (CT_ShapeProperties,
// CT_TableCellProperties) fix the order of their children, so inserting an
// element means finding its position among siblings — and a flat regexp would
// also match nested elements, for example the <a:solidFill> inside an <a:ln>.
//
// It returns nil for anything that is not a plain sequence of elements and
// whitespace, so callers can fall back to leaving the fragment untouched.
func SplitTopLevelXMLElements(fragment string) []XMLElement {
	var elements []XMLElement
	depth := 0
	start := 0
	for i := 0; i < len(fragment); {
		if fragment[i] != '<' {
			if depth == 0 && !isXMLSpace(fragment[i]) {
				return nil
			}
			i++
			continue
		}
		if skipped, isMarkup := skipCommentOrCDATA(fragment[i:]); isMarkup {
			if skipped == 0 {
				return nil
			}
			i += skipped
			continue
		}
		tagEnd := strings.IndexByte(fragment[i:], '>')
		if tagEnd < 0 {
			return nil
		}
		elements, depth, start = consumeXMLTag(
			fragment, elements, fragment[i:i+tagEnd+1], i, tagEnd, depth, start,
		)
		if depth < 0 {
			return nil
		}
		i += tagEnd + 1
	}
	if depth != 0 {
		return nil
	}
	return elements
}

// skipCommentOrCDATA reports how many bytes a comment or CDATA section spans.
// The second result is false when the input starts neither; the length is zero
// when the section is unterminated.
func skipCommentOrCDATA(s string) (int, bool) {
	for prefix, suffix := range map[string]string{"<!--": "-->", "<![CDATA[": "]]>"} {
		if !strings.HasPrefix(s, prefix) {
			continue
		}
		end := strings.Index(s, suffix)
		if end < 0 {
			return 0, true
		}
		return end + len(suffix), true
	}
	return 0, false
}

func consumeXMLTag(
	fragment string,
	elements []XMLElement,
	tag string,
	i, tagEnd, depth, start int,
) ([]XMLElement, int, int) {
	switch {
	case strings.HasPrefix(tag, "</"):
		depth--
		if depth == 0 {
			elements = append(elements, XMLElement{
				Name:  XMLTagName(fragment[start:]),
				Start: start,
				XML:   fragment[start : i+tagEnd+1],
			})
		}
	case strings.HasSuffix(tag, "/>"):
		if depth == 0 {
			elements = append(elements, XMLElement{Name: XMLTagName(tag), Start: i, XML: tag})
		}
	default:
		if depth == 0 {
			start = i
		}
		depth++
	}
	return elements, depth, start
}

// XMLTagName returns the qualified name of an opening tag.
func XMLTagName(tag string) string {
	name := strings.TrimPrefix(tag, "<")
	for i := range len(name) {
		if isXMLSpace(name[i]) || name[i] == '>' || name[i] == '/' {
			return name[:i]
		}
	}
	return name
}

func isXMLSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

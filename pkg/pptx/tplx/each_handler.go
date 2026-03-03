package tplx

import (
	"regexp"
	"strings"
)

// eachSlidePattern detects {{#each KEY}} at the paragraph/shape level of a slide.
var (
	eachSlidePattern = regexp.MustCompile(`\{\{#each\s+([a-zA-Z0-9_]+)\s*\}\}`)
	eachEndPattern   = regexp.MustCompile(`\{\{/each\s*\}\}`)
)

// eachCondition holds parsed each-block metadata.
type eachCondition struct {
	key     string // context key for the slice
	isSlide bool   // true if the entire slide is a loop slide
}

// detectEachSlide checks whether the slide XML is a slide-level loop template.
// A slide is a loop template only when {{#each KEY}} appears OUTSIDE any <a:tbl>
// element (table-row loops are handled separately by expandTableRows).
func detectEachSlide(slideXML []byte) *eachCondition {
	content := string(slideXML)

	// Find position of first {{#each}} token.
	m := eachSlidePattern.FindStringSubmatchIndex(content)
	if m == nil {
		return nil
	}
	tokenStart := m[0]

	// If the token is inside a <a:tbl>…</a:tbl> it is a table-row loop, not a slide loop.
	tblOpen := strings.Index(content, "<a:tbl")
	if tblOpen >= 0 && tblOpen < tokenStart {
		tblClose := strings.Index(content[tblOpen:], "</a:tbl>")
		if tblClose >= 0 && tblOpen+tblClose+len("</a:tbl>") > tokenStart {
			return nil
		}
	}

	key := content[m[2]:m[3]]
	return &eachCondition{key: key, isSlide: true}
}

// expandSlide takes the bytes of a slide-level loop template and the data rows,
// returning one expanded slide XML per row.
// The {{#each KEY}} / {{/each}} markers are removed from each copy.
func expandSlide(slideXML []byte, rows []Row, ctx Context) [][]byte {
	// Strip the #each / /each tokens from the text.
	clean := eachSlidePattern.ReplaceAll(slideXML, nil)
	clean = eachEndPattern.ReplaceAll(clean, nil)

	result := make([][]byte, 0, len(rows))
	for _, row := range rows {
		expanded := interpolateXMLPart(clean, ctx, row, false)
		result = append(result, expanded)
	}
	return result
}

// ── table-row loop ────────────────────────────────────────────────────────────

// expandTableRows scans slide XML for {{#each KEY}} inside a <a:tr> element and
// expands the row N times (once per item in ctx[key]).
// Returns the modified XML with the template row replaced by expanded rows.
//
// Strategy: find the {{#each}} token position, then scan backward to find the
// enclosing <a:tr…> start and forward to find its </a:tr> end.
func expandTableRows(tableXML []byte, ctx Context) []byte {
	content := string(tableXML)

	// Find the {{#each KEY}} token.
	m := eachSlidePattern.FindStringSubmatchIndex(content)
	if m == nil {
		return tableXML
	}
	tokenStart := m[0]
	key := content[m[2]:m[3]]

	// Scan backward from tokenStart to find the opening <a:tr (note no '>' — attribute may follow).
	trOpen := strings.LastIndex(content[:tokenStart], "<a:tr")
	if trOpen < 0 {
		return tableXML
	}

	// Scan forward from trOpen for the matching </a:tr>.
	trCloseRel := strings.Index(content[trOpen:], "</a:tr>")
	if trCloseRel < 0 {
		return tableXML
	}
	trClose := trOpen + trCloseRel + len("</a:tr>")
	rowContent := content[trOpen:trClose]

	// Remove each/end markers from the row template to get the clean row XML.
	cleanRow := eachSlidePattern.ReplaceAllString(rowContent, "")
	cleanRow = eachEndPattern.ReplaceAllString(cleanRow, "")

	// Fetch the rows from context.
	rowsAny, ok := ctx[key]
	if !ok {
		return []byte(content[:trOpen] + content[trClose:])
	}
	rows, ok := toRows(rowsAny)
	if !ok || len(rows) == 0 {
		return []byte(content[:trOpen] + content[trClose:])
	}

	// Build expanded rows.
	var expanded strings.Builder
	for _, row := range rows {
		rowXML := interpolateXMLPart([]byte(cleanRow), ctx, row, false)
		expanded.Write(rowXML)
	}

	return []byte(content[:trOpen] + expanded.String() + content[trClose:])
}

// toRows converts an any value to []Row.
func toRows(v any) ([]Row, bool) {
	switch rv := v.(type) {
	case []Row:
		return rv, true
	case []map[string]string:
		rows := make([]Row, len(rv))
		for i, m := range rv {
			rows[i] = m
		}
		return rows, true
	}
	return nil, false
}

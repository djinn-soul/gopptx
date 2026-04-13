package urlfetch

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

// buildTable converts raw HTML table rows to a table with a bold header row.
func buildTable(rows [][]string) tables.Table {
	if len(rows) == 0 {
		return tables.NewTable(nil)
	}

	cols := 0
	for _, r := range rows {
		if len(r) > cols {
			cols = len(r)
		}
	}

	const totalWidthEMU int64 = 8230200
	colW := totalWidthEMU / int64(cols)
	colWidths := make([]styling.Length, cols)
	for i := range colWidths {
		colWidths[i] = styling.Emu(colW)
	}

	tbl := tables.NewTable(colWidths)
	for i, rawRow := range rows {
		norm := make([]string, cols)
		copy(norm, rawRow)

		if i == 0 {
			cells := make([]tables.TableCell, cols)
			for j, text := range norm {
				cells[j] = tables.NewTableCell(text).WithBold(true)
			}
			tbl = tbl.AddStyledRow(cells)
		} else {
			tbl = tbl.AddRow(norm)
		}
	}
	return tbl
}

// truncateText preserves full extracted text for slide conversion.
// maxLen is intentionally ignored to avoid silent content loss in URL fetch decks.
func truncateText(text string, _ int) string {
	return strings.TrimSpace(text)
}

func splitTextIntoChunks(text string, maxLen int) []string {
	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return nil
	}
	if len(normalized) <= maxLen {
		return []string{normalized}
	}

	words := strings.Fields(normalized)
	if len(words) == 0 {
		return nil
	}

	chunks := make([]string, 0, len(words)/8+1)
	var current strings.Builder
	for _, word := range words {
		if current.Len() == 0 {
			current.WriteString(word)
			continue
		}
		if current.Len()+1+len(word) <= maxLen {
			current.WriteByte(' ')
			current.WriteString(word)
			continue
		}
		chunks = append(chunks, current.String())
		current.Reset()
		current.WriteString(word)
	}
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}
	return chunks
}

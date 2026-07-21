package chart

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/zipfast"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type Kind int

const (
	KindCategory Kind = iota
	KindScatter
	KindBubble
)

const (
	excelDataStartRow      = 2
	scatterHeaderColFactor = 3
)

func GenerateExcelForChart(categories []string, values []float64) ([]byte, error) {
	if len(categories) != len(values) {
		return nil, fmt.Errorf("categories and values length mismatch: %d vs %d", len(categories), len(values))
	}
	seriesName := "Series 1"
	return generateExcelSheetBinary(
		[]string{"Category", seriesName},
		buildCategoryRows(categories, nil, []common.ChartSeriesData{{Values: values}}),
	)
}

func GenerateExcelForChartUpdate(kind Kind, req common.ChartDataUpdate) ([]byte, error) {
	switch kind {
	case KindScatter:
		headers, rows := buildScatterSheet(req.Series, false)
		return generateExcelSheetBinary(headers, rows)
	case KindBubble:
		headers, rows := buildScatterSheet(req.Series, true)
		return generateExcelSheetBinary(headers, rows)
	default:
		headers := categoryHeaders(req)
		for i, s := range req.Series {
			name := fmt.Sprintf("Series %d", i+1)
			if s.Name != nil && strings.TrimSpace(*s.Name) != "" {
				name = *s.Name
			}
			headers = append(headers, name)
		}
		rows := buildCategoryRows(req.Categories, req.MultiLevelCategories, req.Series)
		return generateExcelSheetBinary(headers, rows)
	}
}

func generateExcelSheetBinary(headers []string, rows [][]string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := zipfast.NewWriter(buf)

	if err := writeZipFile(zw, "[Content_Types].xml", ExcelContentTypesXML); err != nil {
		return nil, err
	}
	if err := writeZipFile(zw, "_rels/.rels", ExcelPackageRelsXML); err != nil {
		return nil, err
	}
	if err := writeZipFile(zw, "xl/workbook.xml", ExcelWorkbookXML); err != nil {
		return nil, err
	}
	if err := writeZipFile(zw, "xl/_rels/workbook.xml.rels", ExcelWorkbookRelsXML); err != nil {
		return nil, err
	}
	if err := writeZipFile(zw, "xl/styles.xml", ExcelStylesXML); err != nil {
		return nil, err
	}

	sheetXML, err := generateSheetXML(headers, rows)
	if err != nil {
		return nil, err
	}
	if writeErr := writeZipFile(zw, "xl/worksheets/sheet1.xml", sheetXML); writeErr != nil {
		return nil, writeErr
	}

	if closeErr := zw.Close(); closeErr != nil {
		return nil, closeErr
	}
	return buf.Bytes(), nil
}

func writeZipFile(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}

func generateSheetXML(headers []string, rows [][]string) (string, error) {
	if len(headers) == 0 {
		return "", errors.New("headers cannot be empty")
	}

	xmlRows := fmt.Sprintf(`<row r="1" spans="1:%d">`, len(headers))
	var sbHeaders strings.Builder
	for i, h := range headers {
		cell := ColumnName(i + 1)
		fmt.Fprintf(&sbHeaders, `<c r="%s1" t="inlineStr"><is><t>%s</t></is></c>`, cell, simpleXMLEscape(h))
	}
	xmlRows += sbHeaders.String()
	xmlRows += `</row>`

	var sbRows strings.Builder
	for i, row := range rows {
		rowNum := i + excelDataStartRow
		fmt.Fprintf(&sbRows, `<row r="%d" spans="1:%d">`, rowNum, len(headers))
		var sbCols strings.Builder
		for col := range headers {
			val := ""
			if col < len(row) {
				val = row[col]
			}
			cell := ColumnName(col + 1)
			if isNumberLiteral(val) {
				fmt.Fprintf(&sbCols, `<c r="%s%d"><v>%s</v></c>`, cell, rowNum, val)
			} else {
				fmt.Fprintf(
					&sbCols,
					`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`,
					cell,
					rowNum,
					simpleXMLEscape(val),
				)
			}
		}
		sbRows.WriteString(sbCols.String())
		sbRows.WriteString(`</row>`)
	}
	xmlRows += sbRows.String()

	return fmt.Sprintf(ExcelSheetTemplate, xmlRows), nil
}

func simpleXMLEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

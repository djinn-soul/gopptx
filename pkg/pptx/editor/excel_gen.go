package editor

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	excelDataStartRow      = 2
	scatterHeaderColFactor = 3
	excelColumnBase        = 26
)

// generateExcelForChart creates a minimal .xlsx file content (as []byte)
// suitable for a PowerPoint chart data source.
//
// It generates a single sheet "Sheet1" with:
// - Header row: "Category", "Series 1"
// - Data rows: category[i], value[i]
//
// Limitations:
// - Single series only (for now)
// - No styling beyond basics
// - No shared strings optimization (inline strings for simplicity).
func generateExcelForChart(categories []string, values []float64) ([]byte, error) {
	if len(categories) != len(values) {
		return nil, fmt.Errorf("categories and values length mismatch: %d vs %d", len(categories), len(values))
	}
	seriesName := "Series 1"
	return generateExcelSheetBinary(
		[]string{"Category", seriesName},
		buildCategoryRows(categories, []common.ChartSeriesData{{Values: values}}),
	)
}

func generateExcelForChartUpdate(kind chartKind, req common.ChartDataUpdate) ([]byte, error) {
	switch kind {
	case chartKindScatter:
		headers, rows := buildScatterSheet(req.Series, false)
		return generateExcelSheetBinary(headers, rows)
	case chartKindBubble:
		headers, rows := buildScatterSheet(req.Series, true)
		return generateExcelSheetBinary(headers, rows)
	default:
		headers := []string{"Category"}
		for i, s := range req.Series {
			name := fmt.Sprintf("Series %d", i+1)
			if s.Name != nil && strings.TrimSpace(*s.Name) != "" {
				name = *s.Name
			}
			headers = append(headers, name)
		}
		rows := buildCategoryRows(req.Categories, req.Series)
		return generateExcelSheetBinary(headers, rows)
	}
}

func generateExcelSheetBinary(headers []string, rows [][]string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// 1. [Content_Types].xml
	if err := writeZipFile(zw, "[Content_Types].xml", ExcelContentTypesXML); err != nil {
		return nil, err
	}

	// 2. _rels/.rels
	if err := writeZipFile(zw, "_rels/.rels", ExcelPackageRelsXML); err != nil {
		return nil, err
	}

	// 3. xl/workbook.xml
	if err := writeZipFile(zw, "xl/workbook.xml", ExcelWorkbookXML); err != nil {
		return nil, err
	}

	// 4. xl/_rels/workbook.xml.rels
	if err := writeZipFile(zw, "xl/_rels/workbook.xml.rels", ExcelWorkbookRelsXML); err != nil {
		return nil, err
	}

	// 5. xl/styles.xml (Minimal)
	if err := writeZipFile(zw, "xl/styles.xml", ExcelStylesXML); err != nil {
		return nil, err
	}

	// 6. xl/worksheets/sheet1.xml (The actual data)
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
	var xmlRowsSb115 strings.Builder
	for i, h := range headers {
		cell := columnName(i + 1)
		xmlRowsSb115.WriteString(
			fmt.Sprintf(`<c r="%s1" t="inlineStr"><is><t>%s</t></is></c>`, cell, simpleXMLEscape(h)),
		)
	}
	xmlRows += xmlRowsSb115.String()
	xmlRows += `</row>`

	var xmlRowsSb121 strings.Builder
	for i, row := range rows {
		rowNum := i + excelDataStartRow
		xmlRowsSb121.WriteString(fmt.Sprintf(`<row r="%d" spans="1:%d">`, rowNum, len(headers)))
		var xmlRowsSb124 strings.Builder
		for col := range headers {
			val := ""
			if col < len(row) {
				val = row[col]
			}
			cell := columnName(col + 1)
			if isNumberLiteral(val) {
				xmlRowsSb124.WriteString(fmt.Sprintf(`<c r="%s%d"><v>%s</v></c>`, cell, rowNum, val))
			} else {
				xmlRowsSb124.WriteString(
					fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`, cell, rowNum, simpleXMLEscape(val)),
				)
			}
		}
		xmlRowsSb121.WriteString(xmlRowsSb124.String())
		xmlRowsSb121.WriteString(`</row>`)
	}
	xmlRows += xmlRowsSb121.String()

	return fmt.Sprintf(ExcelSheetTemplate, xmlRows), nil
}

func simpleXMLEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func buildCategoryRows(categories []string, series []common.ChartSeriesData) [][]string {
	rowCount := len(categories)
	if rowCount == 0 && len(series) > 0 && len(series[0].Categories) > 0 {
		rowCount = len(series[0].Categories)
	}
	rows := make([][]string, rowCount)
	for i := range rowCount {
		row := make([]string, 0, 1+len(series))
		cat := ""
		if len(categories) > i {
			cat = categories[i]
		} else if len(series) > 0 && len(series[0].Categories) > i {
			cat = series[0].Categories[i]
		}
		row = append(row, cat)
		for _, s := range series {
			if len(s.Values) > i {
				row = append(row, strconv.FormatFloat(s.Values[i], 'f', -1, 64))
			} else {
				row = append(row, "")
			}
		}
		rows[i] = row
	}
	return rows
}

func buildScatterSheet(series []common.ChartSeriesData, withSizes bool) ([]string, [][]string) {
	headers := make([]string, 0, len(series)*scatterHeaderColFactor)
	maxRows := 0
	for i, s := range series {
		headers = append(headers, fmt.Sprintf("X%d", i+1), fmt.Sprintf("Y%d", i+1))
		if withSizes {
			headers = append(headers, fmt.Sprintf("S%d", i+1))
		}
		if len(s.XValues) > maxRows {
			maxRows = len(s.XValues)
		}
	}
	rows := make([][]string, maxRows)
	for r := range maxRows {
		row := make([]string, 0, len(headers))
		for _, s := range series {
			if len(s.XValues) > r {
				row = append(row, strconv.FormatFloat(s.XValues[r], 'f', -1, 64))
			} else {
				row = append(row, "")
			}
			if len(s.YValues) > r {
				row = append(row, strconv.FormatFloat(s.YValues[r], 'f', -1, 64))
			} else {
				row = append(row, "")
			}
			if withSizes {
				if len(s.Sizes) > r {
					row = append(row, strconv.FormatFloat(s.Sizes[r], 'f', -1, 64))
				} else {
					row = append(row, "")
				}
			}
		}
		rows[r] = row
	}
	return headers, rows
}

func columnName(n int) string {
	name := ""
	for n > 0 {
		n--
		name = string(rune('A'+(n%excelColumnBase))) + name
		n /= excelColumnBase
	}
	return name
}

func isNumberLiteral(s string) bool {
	if s == "" {
		return false
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return false
	}
	return true
}

// ExcelContentTypesXML and related constants form a minimal valid Excel package.
const ExcelContentTypesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/><Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/></Types>`

const ExcelPackageRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`

const ExcelWorkbookXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><bookViews><workbookView xWindow="0" yWindow="0" windowWidth="20490" windowHeight="13245"/></bookViews><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>`

const ExcelWorkbookRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/></Relationships>`

const ExcelStylesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fonts count="1"><font><sz val="11"/><color theme="1"/><name val="Calibri"/><family val="2"/><scheme val="minor"/></font></fonts><fills count="2"><fill><patternFill patternType="none"/></fill><fill><patternFill patternType="gray125"/></fill></fills><borders count="1"><border><left/><right/><top/><bottom/><diagonal/></border></borders><cellStyleXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellStyleXfs><cellXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0"/></cellXfs><cellStyles count="1"><cellStyle name="Normal" xfId="0" builtinId="0"/></cellStyles><dxfs count="0"/><tableStyles count="0" defaultTableStyle="TableStyleMedium2" defaultPivotStyle="PivotStyleLight16"/></styleSheet>`

const ExcelSheetTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheetData>%s</sheetData></worksheet>`

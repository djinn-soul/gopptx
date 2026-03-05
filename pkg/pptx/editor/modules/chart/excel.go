package chart

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	excelColumnBase        = 26
)

func GenerateExcelForChart(categories []string, values []float64) ([]byte, error) {
	if len(categories) != len(values) {
		return nil, fmt.Errorf("categories and values length mismatch: %d vs %d", len(categories), len(values))
	}
	seriesName := "Series 1"
	return generateExcelSheetBinary(
		[]string{"Category", seriesName},
		buildCategoryRows(categories, []common.ChartSeriesData{{Values: values}}),
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
		sbHeaders.WriteString(
			fmt.Sprintf(`<c r="%s1" t="inlineStr"><is><t>%s</t></is></c>`, cell, simpleXMLEscape(h)),
		)
	}
	xmlRows += sbHeaders.String()
	xmlRows += `</row>`

	var sbRows strings.Builder
	for i, row := range rows {
		rowNum := i + excelDataStartRow
		sbRows.WriteString(fmt.Sprintf(`<row r="%d" spans="1:%d">`, rowNum, len(headers)))
		var sbCols strings.Builder
		for col := range headers {
			val := ""
			if col < len(row) {
				val = row[col]
			}
			cell := ColumnName(col + 1)
			if isNumberLiteral(val) {
				sbCols.WriteString(fmt.Sprintf(`<c r="%s%d"><v>%s</v></c>`, cell, rowNum, val))
			} else {
				sbCols.WriteString(
					fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`, cell, rowNum, simpleXMLEscape(val)),
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
	headers, maxRows := scatterHeadersAndMaxRows(series, withSizes)
	rows := make([][]string, maxRows)
	for rowIdx := range maxRows {
		rows[rowIdx] = scatterRowValues(series, rowIdx, withSizes, len(headers))
	}
	return headers, rows
}

func scatterHeadersAndMaxRows(series []common.ChartSeriesData, withSizes bool) ([]string, int) {
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
	return headers, maxRows
}

func scatterRowValues(series []common.ChartSeriesData, rowIdx int, withSizes bool, capacity int) []string {
	row := make([]string, 0, capacity)
	for _, s := range series {
		row = append(row, scatterValueAt(s.XValues, rowIdx))
		row = append(row, scatterValueAt(s.YValues, rowIdx))
		if withSizes {
			row = append(row, scatterValueAt(s.Sizes, rowIdx))
		}
	}
	return row
}

func scatterValueAt(values []float64, idx int) string {
	if len(values) <= idx {
		return ""
	}
	return strconv.FormatFloat(values[idx], 'f', -1, 64)
}

func ColumnName(n int) string {
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

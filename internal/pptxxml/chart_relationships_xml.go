package pptxxml

// ChartEmbeddingContentType is the content type of the workbook a chart is
// built from.
const ChartEmbeddingContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// chartEmbeddingRelType is the OPC relationship a chart uses to reach its
// workbook: the whole package is embedded, so it is the package reltype rather
// than a document-level one.
const chartEmbeddingRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package"

// ChartRelationships renders ppt/charts/_rels/chartN.xml.rels, which binds a
// chart part to the workbook holding its data.
func ChartRelationships(relID, target string) string {
	return xmlHeader + `
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="` + Escape(relID) + `" Type="` + chartEmbeddingRelType + `" Target="` + Escape(target) + `"/>
</Relationships>`
}

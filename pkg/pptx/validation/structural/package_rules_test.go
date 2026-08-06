package structural

import "testing"

func issueCodes(issues []Issue) map[IssueCode]int {
	counts := make(map[IssueCode]int, len(issues))
	for _, issue := range issues {
		counts[issue.Code]++
	}
	return counts
}

// The validator used to pass a deck whose chart had no workbook at all, which
// is exactly the state "Edit Data" cannot open.
func TestValidatorReportsChartWithoutWorkbook(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/charts/chart1.xml": []byte(`<c:chartSpace xmlns:c="c"><c:chart/></c:chartSpace>`),
	}}
	issues := NewValidator(m).Validate()

	if issueCodes(issues)[CodeIncompleteChartPackage] == 0 {
		t.Fatalf("no chart-package issue reported: %v", issues)
	}
}

// A chart that names external data but ships no relationships part cannot
// resolve it.
func TestValidatorReportsChartWithoutRelationships(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/charts/chart1.xml": []byte(
			`<c:chartSpace xmlns:c="c"><c:chart/><c:externalData r:id="rId1"/></c:chartSpace>`,
		),
	}}
	issues := NewValidator(m).Validate()

	found := false
	for _, issue := range issues {
		if issue.Code == CodeIncompleteChartPackage && issue.Severity == SeverityError {
			found = true
		}
	}
	if !found {
		t.Fatalf("no error for a chart with external data and no rels: %v", issues)
	}
}

// The handout master defect: part present, never declared.
func TestValidatorReportsUndeclaredHandoutMaster(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/presentation.xml":                  []byte(`<p:presentation><p:sldIdLst/></p:presentation>`),
		"ppt/handoutMasters/handoutMaster1.xml": []byte(`<p:handoutMaster/>`),
	}}
	issues := NewValidator(m).Validate()

	if issueCodes(issues)[CodeUndeclaredPart] == 0 {
		t.Fatalf("undeclared handout master not reported: %v", issues)
	}
}

func TestValidatorReportsDuplicateMasterIDs(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/presentation.xml": []byte(
			`<p:presentation><p:sldMasterIdLst>` +
				`<p:sldMasterId id="2147483648" r:id="rId1"/>` +
				`<p:sldMasterId id="2147483648" r:id="rId2"/>` +
				`</p:sldMasterIdLst></p:presentation>`,
		),
	}}
	issues := NewValidator(m).Validate()

	if issueCodes(issues)[CodeDuplicateID] == 0 {
		t.Fatalf("duplicate master id not reported: %v", issues)
	}
}

// A master that reaches no theme and no layout is unusable.
func TestValidatorReportsIncompleteMaster(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/slideMasters/slideMaster1.xml": []byte(`<p:sldMaster/>`),
		"ppt/slideMasters/_rels/slideMaster1.xml.rels": []byte(
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`,
		),
	}}
	issues := NewValidator(m).Validate()

	if issueCodes(issues)[CodeIncompleteMaster] < 2 {
		t.Fatalf("expected both theme and layout issues: %v", issues)
	}
}

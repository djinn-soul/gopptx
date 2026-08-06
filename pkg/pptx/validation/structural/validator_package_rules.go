package structural

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"
)

// The rules below are the ones that would have caught the package defects a
// shallow check misses: a chart with no workbook, a master with no theme, a
// handout written but never declared. Each is a separate family so a report
// says which one failed.

const (
	// CodeIncompleteChartPackage marks a chart part whose workbook or
	// relationships are missing.
	CodeIncompleteChartPackage IssueCode = "INCOMPLETE_CHART_PACKAGE"
	// CodeIncompleteMaster marks a slide master that does not reach a theme or
	// any layout.
	CodeIncompleteMaster IssueCode = "INCOMPLETE_MASTER"
	// CodeUndeclaredPart marks a part that is present and related but never
	// named in presentation.xml.
	CodeUndeclaredPart IssueCode = "UNDECLARED_PART"
	// CodeDuplicateID marks a repeated master or layout id.
	CodeDuplicateID IssueCode = "DUPLICATE_ID"
)

// Package paths named in more than one rule.
const (
	firstMasterPartPath = "ppt/slideMasters/slideMaster1.xml"
	firstLayoutPartPath = "ppt/slideLayouts/slideLayout1.xml"
	corePropsPartPath   = "docProps/core.xml"
	appPropsPartPath    = "docProps/app.xml"
)

//nolint:gochecknoglobals // read-only lookup tables, never mutated
var (
	// standardRequiredParts are the parts PowerPoint writes into every deck and
	// expects to find when opening one. Missing any of them does not always stop
	// a deck opening, so they are warnings rather than errors — unlike the four
	// in requiredParts, without which the package is not a package.
	standardRequiredParts = map[string]string{
		"ppt/presProps.xml":    "Presentation properties",
		"ppt/viewProps.xml":    "View properties",
		"ppt/tableStyles.xml":  "Table styles",
		"ppt/theme/theme1.xml": "Theme",
		firstMasterPartPath:    "Slide master",
		firstLayoutPartPath:    "Slide layout",
		corePropsPartPath:      "Core properties",
		appPropsPartPath:       "Extended properties",
	}

	chartPartPattern      = regexp.MustCompile(`^ppt/charts/chart\d+\.xml$`)
	notesSlidePattern     = regexp.MustCompile(`^ppt/notesSlides/notesSlide\d+\.xml$`)
	masterPartPattern     = regexp.MustCompile(`^ppt/slideMasters/slideMaster\d+\.xml$`)
	externalDataPattern   = regexp.MustCompile(`<c:externalData\b`)
	masterIDPattern       = regexp.MustCompile(`<p:sldMasterId\s+id="(\d+)"`)
	layoutIDPattern       = regexp.MustCompile(`<p:sldLayoutId\s+id="(\d+)"`)
	handoutMasterPartPath = "ppt/handoutMasters/handoutMaster1.xml"
	notesMasterPartPath   = "ppt/notesMasters/notesMaster1.xml"
	themeRelType          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	layoutRelType         = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
)

// checkStandardParts reports the parts PowerPoint expects in every deck.
func (v *Validator) checkStandardParts() {
	for p, desc := range standardRequiredParts {
		if !v.provider.Has(p) {
			v.issues = append(v.issues, Issue{
				Code:        CodeMissingPart,
				Severity:    SeverityWarning,
				Path:        p,
				Description: "Missing part PowerPoint expects: " + desc,
			})
		}
	}
}

// checkChartPackages reports charts that ship without the workbook their data
// comes from. This is the rule that catches a chart PowerPoint can draw but
// cannot open the data for.
func (v *Validator) checkChartPackages() {
	for _, p := range v.keys {
		if !chartPartPattern.MatchString(p) {
			continue
		}
		data, ok := v.provider.Get(p)
		if !ok {
			continue
		}
		if !externalDataPattern.Match(data) {
			v.issues = append(v.issues, Issue{
				Code:        CodeIncompleteChartPackage,
				Severity:    SeverityWarning,
				Path:        p,
				Description: "Chart has no <c:externalData>, so its data cannot be edited in PowerPoint",
			})
			continue
		}
		relsPath := path.Join(path.Dir(p), "_rels", path.Base(p)+".rels")
		if !v.provider.Has(relsPath) {
			v.issues = append(v.issues, Issue{
				Code:        CodeIncompleteChartPackage,
				Severity:    SeverityError,
				Path:        p,
				Description: "Chart references external data but has no relationships part",
			})
		}
	}
}

// checkSlideMasters reports masters that reach no theme or no layout.
// A package with no master at all is reported once, by checkStandardParts.
func (v *Validator) checkSlideMasters() {
	for _, p := range v.keys {
		if !masterPartPattern.MatchString(p) {
			continue
		}
		relsPath := path.Join(path.Dir(p), "_rels", path.Base(p)+".rels")
		data, ok := v.provider.Get(relsPath)
		if !ok {
			v.issues = append(v.issues, Issue{
				Code:        CodeIncompleteMaster,
				Severity:    SeverityError,
				Path:        p,
				Description: "Slide master has no relationships part",
			})
			continue
		}
		hasTheme, layouts := masterRelationshipSummary(data)
		if !hasTheme {
			v.issues = append(v.issues, Issue{
				Code:        CodeIncompleteMaster,
				Severity:    SeverityError,
				Path:        p,
				Description: "Slide master does not reference a theme",
			})
		}
		if layouts == 0 {
			v.issues = append(v.issues, Issue{
				Code:        CodeIncompleteMaster,
				Severity:    SeverityError,
				Path:        p,
				Description: "Slide master references no slide layout",
			})
		}
	}
}

// masterRelationshipSummary reports whether the master reaches a theme, and how
// many layouts it reaches.
func masterRelationshipSummary(data []byte) (bool, int) {
	var rels relationshipsXML
	if err := xml.NewDecoder(bytes.NewReader(data)).Decode(&rels); err != nil {
		return false, 0
	}
	hasTheme := false
	layoutCount := 0
	for _, rel := range rels.Relationships {
		switch rel.Type {
		case themeRelType:
			hasTheme = true
		case layoutRelType:
			layoutCount++
		}
	}
	return hasTheme, layoutCount
}

// checkNotesAndHandoutMasters reports the two masters that can be present in
// the package yet unreachable from presentation.xml.
func (v *Validator) checkNotesAndHandoutMasters() {
	presentation, hasPresentation := v.provider.Get(presentationPartPath)
	if !hasPresentation {
		return
	}

	if v.provider.Has(handoutMasterPartPath) &&
		!bytes.Contains(presentation, []byte("<p:handoutMasterIdLst>")) {
		v.issues = append(v.issues, Issue{
			Code:        CodeUndeclaredPart,
			Severity:    SeverityError,
			Path:        handoutMasterPartPath,
			Description: "Handout master is present but presentation.xml has no handoutMasterIdLst",
		})
	}

	if !slices.ContainsFunc(v.keys, notesSlidePattern.MatchString) {
		return
	}
	if !v.provider.Has(notesMasterPartPath) {
		v.issues = append(v.issues, Issue{
			Code:        CodeMissingPart,
			Severity:    SeverityError,
			Path:        notesMasterPartPath,
			Description: "Package has notes slides but no notes master",
		})
		return
	}
	if !bytes.Contains(presentation, []byte("<p:notesMasterIdLst>")) {
		v.issues = append(v.issues, Issue{
			Code:        CodeUndeclaredPart,
			Severity:    SeverityError,
			Path:        notesMasterPartPath,
			Description: "Notes master is present but presentation.xml has no notesMasterIdLst",
		})
	}
}

// checkMasterIDUniqueness reports repeated master or layout ids, which make
// PowerPoint resolve the wrong one.
func (v *Validator) checkMasterIDUniqueness() {
	presentation, ok := v.provider.Get(presentationPartPath)
	if ok {
		v.reportDuplicateIDs(presentationPartPath, "slide master", masterIDPattern.FindAllSubmatch(presentation, -1))
	}
	for _, p := range v.keys {
		if !masterPartPattern.MatchString(p) {
			continue
		}
		data, found := v.provider.Get(p)
		if !found {
			continue
		}
		v.reportDuplicateIDs(p, "slide layout", layoutIDPattern.FindAllSubmatch(data, -1))
	}
}

func (v *Validator) reportDuplicateIDs(partPath, kind string, matches [][][]byte) {
	seen := make(map[string]bool, len(matches))
	for _, match := range matches {
		id := string(match[1])
		if seen[id] {
			v.issues = append(v.issues, Issue{
				Code:        CodeDuplicateID,
				Severity:    SeverityError,
				Path:        partPath,
				Description: fmt.Sprintf("Duplicate %s id %s", kind, id),
			})
			continue
		}
		seen[id] = true
	}
}

// checkThemePart reports a theme part that is not a theme.
func (v *Validator) checkThemePart() {
	for _, p := range v.keys {
		if !strings.HasPrefix(p, "ppt/theme/") || !strings.HasSuffix(p, ".xml") {
			continue
		}
		data, ok := v.provider.Get(p)
		if !ok {
			continue
		}
		if !bytes.Contains(data, []byte("<a:theme")) {
			v.issues = append(v.issues, Issue{
				Code:        CodeInvalidXML,
				Severity:    SeverityError,
				Path:        p,
				Description: "Theme part has no a:theme root element",
			})
		}
	}
}

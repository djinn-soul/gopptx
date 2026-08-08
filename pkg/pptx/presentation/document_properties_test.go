package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// common.CoreProperties declares fourteen fields; the generator carried four of
// them to the file and hardcoded the rest, so everything else a caller set was
// dropped.
func TestGeneratorWritesAllCoreProperties(t *testing.T) {
	meta := Metadata{Metadata: common.Metadata{
		Title: "Deck",
		CoreProperties: common.CoreProperties{
			Keywords:       "parity,ooxml",
			Category:       "Reports",
			ContentStatus:  "Final",
			Identifier:     "DECK-1",
			Language:       "en-GB",
			Version:        "3.2",
			LastModifiedBy: "Reviewer",
			Revision:       "7",
			Created:        "2020-01-02T03:04:05Z",
			Modified:       "2021-02-03T04:05:06Z",
			LastPrinted:    "2021-03-04T05:06:07Z",
		},
	}}
	parts := buildPackageParts(t, meta, []elements.SlideContent{elements.NewSlide("S1")})
	core := parts["docProps/core.xml"]

	wants := []string{
		"<cp:keywords>parity,ooxml</cp:keywords>",
		"<cp:category>Reports</cp:category>",
		"<cp:contentStatus>Final</cp:contentStatus>",
		"<dc:identifier>DECK-1</dc:identifier>",
		"<dc:language>en-GB</dc:language>",
		"<cp:version>3.2</cp:version>",
		"<cp:lastModifiedBy>Reviewer</cp:lastModifiedBy>",
		"<cp:revision>7</cp:revision>",
		"<cp:lastPrinted>2021-03-04T05:06:07Z</cp:lastPrinted>",
		">2020-01-02T03:04:05Z</dcterms:created>",
		">2021-02-03T04:05:06Z</dcterms:modified>",
	}
	for _, want := range wants {
		if !strings.Contains(core, want) {
			t.Errorf("core.xml is missing %s:\n%s", want, core)
		}
	}
}

// app.xml hardcoded the application, emitted no Company or Manager, and wrote
// HiddenSlides as zero even for a deck with hidden slides.
func TestGeneratorWritesAppProperties(t *testing.T) {
	meta := Metadata{Metadata: common.Metadata{
		AppProperties: common.AppProperties{
			Application: "Acme Deck Builder",
			AppVersion:  "2.5000",
			Company:     "Acme",
			Manager:     "R. Runner",
		},
	}}
	hidden := elements.NewSlide("Hidden")
	hidden.Hidden = true
	parts := buildPackageParts(t, meta, []elements.SlideContent{elements.NewSlide("S1"), hidden})
	app := parts["docProps/app.xml"]

	wants := []string{
		"<Application>Acme Deck Builder</Application>",
		"<AppVersion>2.5000</AppVersion>",
		"<Company>Acme</Company>",
		"<Manager>R. Runner</Manager>",
		"<HiddenSlides>1</HiddenSlides>",
	}
	for _, want := range wants {
		if !strings.Contains(app, want) {
			t.Errorf("app.xml is missing %s:\n%s", want, app)
		}
	}
}

package editor

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Count elements of docProps/app.xml, compiled once: renderAppProperties runs
// on every save.
//
// "Slides" and "HiddenSlides" cannot cross-match even though one name ends with
// the other: "<Slides>" does not occur inside "<HiddenSlides>", because the
// character before "Slides" there is "n", not "<".
//
//nolint:gochecknoglobals // compiled once; renderAppProperties runs on every save
var (
	appSlidesPattern       = appPropertyPattern("Slides")
	appNotesPattern        = appPropertyPattern("Notes")
	appHiddenSlidesPattern = appPropertyPattern("HiddenSlides")
	appPropertiesClose     = regexp.MustCompile(`</(?:[A-Za-z_][\w.-]*:)?Properties\s*>`)
)

func appPropertyPattern(tag string) *regexp.Regexp {
	qualified := `(?:[A-Za-z_][\w.-]*:)?` + regexp.QuoteMeta(tag)
	return regexp.MustCompile(`(?s)<` + qualified + `\b[^>]*>.*?</` + qualified + `\s*>`)
}

// appPropertyCounts are the docProps/app.xml values that go stale as soon as a
// slide is added, removed or hidden through the editor. Gmail, Outlook and the
// Windows shell preview a deck from these numbers rather than from the slide
// list, so a stale <Slides> shows the wrong deck.
type appPropertyCounts struct {
	slides       int
	notes        int
	hiddenSlides int
}

func (e *PresentationEditor) appPropertyCounts() appPropertyCounts {
	counts := appPropertyCounts{slides: len(e.slides)}
	// A notes part can be shared by two slide refs after a slide copy, and the
	// inventory can still hold entries for slides that were since removed, so
	// count distinct parts reachable from the current slide list.
	seen := make(map[string]struct{}, len(e.notesInventory))
	for _, slideRef := range e.slides {
		if slideRef.Hidden {
			counts.hiddenSlides++
		}
		notesPart := strings.TrimSpace(e.notesInventory[slideRef.Part])
		if notesPart == "" {
			continue
		}
		if _, dup := seen[notesPart]; dup {
			continue
		}
		seen[notesPart] = struct{}{}
		counts.notes++
	}
	return counts
}

// renderAppProperties returns docProps/app.xml with the slide, notes and hidden
// slide counts brought up to date.
//
// An existing part is patched rather than regenerated: PowerPoint writes
// Company, Manager, TitlesOfParts and the fonts it used into app.xml, and
// replacing the part wholesale would drop all of it. Only a deck that has no
// app.xml at all gets a generated one.
//
// TitlesOfParts is deliberately left as found. It is a vector of per-slide
// titles that only PowerPoint maintains; rewriting it from the editor's slide
// list would also have to rewrite the HeadingPairs entries for fonts and theme
// that sit in the same vector, and nothing reads it for previews.
func (e *PresentationEditor) renderAppProperties() []byte {
	counts := e.appPropertyCounts()

	appXML, ok := e.parts.Get(common.AppPropsPath)
	if !ok || len(bytes.TrimSpace(appXML)) == 0 {
		appXML = []byte(pptxxml.AppProperties(pptxxml.AppPropertiesInfo{
			SlideCount:   counts.slides,
			NotesCount:   counts.notes,
			HiddenSlides: counts.hiddenSlides,
			Width:        e.metadata.SlideSize.Width,
			Height:       e.metadata.SlideSize.Height,
			Application:  e.metadata.AppProperties.Application,
			AppVersion:   e.metadata.AppProperties.AppVersion,
			Company:      e.metadata.AppProperties.Company,
			Manager:      e.metadata.AppProperties.Manager,
		}))
	}

	patched := string(appXML)
	patched = setAppPropertyInt(patched, appSlidesPattern, "Slides", counts.slides)
	patched = setAppPropertyInt(patched, appNotesPattern, "Notes", counts.notes)
	patched = setAppPropertyInt(patched, appHiddenSlidesPattern, "HiddenSlides", counts.hiddenSlides)
	return []byte(patched)
}

// setAppPropertyInt replaces the value of an app.xml count element, appending
// the element if the source document omits it.
func setAppPropertyInt(appXML string, pattern *regexp.Regexp, tag string, value int) string {
	valueText := strconv.Itoa(value)
	if pattern.MatchString(appXML) {
		return pattern.ReplaceAllStringFunc(appXML, func(element string) string {
			openTagEnd := strings.Index(element, ">")
			closeStart := strings.LastIndex(element, "</")
			if openTagEnd < 0 || closeStart < openTagEnd {
				return element
			}
			return element[:openTagEnd+1] + valueText + element[closeStart:]
		})
	}
	closeLocation := appPropertiesClose.FindAllStringIndex(appXML, -1)
	if len(closeLocation) == 0 {
		return appXML
	}
	last := closeLocation[len(closeLocation)-1]
	closing := appXML[last[0]:last[1]]
	prefix := strings.TrimSuffix(strings.TrimPrefix(closing, "</"), "Properties>")
	element := "<" + prefix + tag + ">" + valueText + "</" + prefix + tag + ">\n"
	return appXML[:last[0]] + element + appXML[last[0]:]
}

// ensureAppPropertiesRelationship adds the package-level relationship pointing
// at docProps/app.xml when the source deck has none. A deck that was written
// without app.xml still needs the relationship, or the part is orphaned and
// previewers ignore it.
func ensureAppPropertiesRelationship(rels []common.EditorRelationship) ([]common.EditorRelationship, bool) {
	for _, rel := range rels {
		if rel.Type == common.RelTypeExtendedProps {
			return rels, false
		}
	}
	return append(rels, common.EditorRelationship{
		ID:     "rId" + strconv.Itoa(common.NextRelationshipNumber(rels)),
		Type:   common.RelTypeExtendedProps,
		Target: common.AppPropsPath,
	}), true
}

package export

import (
	"sync"

	"github.com/signintech/gopdf"
)

// Font state belongs to one PDF document, not to the process. Aliases are
// registered on a particular *gopdf.GoPdf, and the metrics behind them describe
// the files that document embedded, so two exports running at once must not see
// each other's: the second would resolve a run to an alias the first registered
// and gopdf would refuse to set it, silently drawing the text in the wrong face.
//
// Each document therefore gets its own documentFonts, keyed by the pdf it
// belongs to and dropped when the export finishes. The index of fonts installed
// on the host stays process-wide (see pdf_font_named.go): that describes the
// machine, and every document sees the same one.

// documentFonts is the font registry of a single PDF document.
type documentFonts struct {
	mu sync.RWMutex
	// metrics maps a registered alias to the vertical metrics of the file
	// embedded under it.
	metrics map[string]ttfLineMetrics
	// named maps a normalized family name to the alias it was embedded under.
	// An empty value records a family this host does not have, so it is looked
	// up only once per document.
	named map[string]string

	sansAlias  string
	serifAlias string
	monoAlias  string
	cjkAlias   string
}

func newDocumentFonts() *documentFonts {
	return &documentFonts{
		metrics:    map[string]ttfLineMetrics{},
		named:      map[string]string{},
		sansAlias:  fontFamilySans,
		serifAlias: fontFamilySans,
		monoAlias:  fontFamilySans,
		cjkAlias:   fontFamilySans,
	}
}

//nolint:gochecknoglobals // One entry per in-flight document; entries are removed when it is written.
var pdfDocumentFonts sync.Map // *gopdf.GoPdf -> *documentFonts

// documentFontsFor returns the font state of pdf, creating it on first use. A
// nil pdf gets a throwaway registry so that helpers stay usable off a document.
func documentFontsFor(pdf *gopdf.GoPdf) *documentFonts {
	if pdf == nil {
		return newDocumentFonts()
	}
	if state, ok := pdfDocumentFonts.Load(pdf); ok {
		fonts, _ := state.(*documentFonts)
		return fonts
	}
	state, _ := pdfDocumentFonts.LoadOrStore(pdf, newDocumentFonts())
	fonts, _ := state.(*documentFonts)
	return fonts
}

// startDocumentFonts gives pdf a clean registry, discarding anything a previous
// export of the same (reused) document left behind.
func startDocumentFonts(pdf *gopdf.GoPdf) {
	if pdf == nil {
		return
	}
	pdfDocumentFonts.Store(pdf, newDocumentFonts())
}

// releaseDocumentFonts drops the registry once the document is written, so a
// long-lived process does not accumulate one per export.
func releaseDocumentFonts(pdf *gopdf.GoPdf) {
	if pdf == nil {
		return
	}
	pdfDocumentFonts.Delete(pdf)
}

func (d *documentFonts) putMetrics(alias string, metrics ttfLineMetrics) {
	d.mu.Lock()
	d.metrics[alias] = metrics
	d.mu.Unlock()
}

func (d *documentFonts) lookupMetrics(alias string) (ttfLineMetrics, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	metrics, ok := d.metrics[alias]
	return metrics, ok
}

func (d *documentFonts) namedAlias(family string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	alias, ok := d.named[family]
	return alias, ok
}

func (d *documentFonts) putNamedAlias(family, alias string) {
	d.mu.Lock()
	d.named[family] = alias
	d.mu.Unlock()
}

func (d *documentFonts) setGenericAliases(sansAlias, serifAlias, monoAlias string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sansAlias = fallbackAlias(sansAlias, fontFamilySans)
	d.serifAlias = fallbackAlias(serifAlias, d.sansAlias)
	d.monoAlias = fallbackAlias(monoAlias, d.sansAlias)
	d.cjkAlias = fallbackAlias(d.cjkAlias, d.sansAlias)
}

func (d *documentFonts) setCJKAlias(alias string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cjkAlias = fallbackAlias(alias, d.sansAlias)
}

// genericAliases returns the sans, serif, mono and CJK aliases of the document.
func (d *documentFonts) genericAliases() (string, string, string, string) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.sansAlias, d.serifAlias, d.monoAlias, d.cjkAlias
}

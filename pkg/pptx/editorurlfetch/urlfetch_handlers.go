// Package editorurlfetch registers the URL-fetch operation handler into the
// editor bridge. It lives in a separate package to break the import cycle:
//
//	editor → urlfetch → pptx → editor
//
// Call Register() once at startup (e.g. from bindings/c/bridge.go init).
package editorurlfetch

import (
	"encoding/json"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"
)

// Register wires the URL-fetch handler into the editor command dispatcher.
// It must be called before any bridge requests are processed.
func Register() {
	editor.RegisterURLFetchHandler(editor.OpURLFetchToSlides, handleURLFetchToSlides)
}

func handleURLFetchToSlides(e *editor.PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := editor.ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := editor.NewPayloadValidator()
	url, ok := v.RequireString(p, "url")
	if !ok {
		return nil, v.Error()
	}

	cfg := urlfetch.DefaultConfig()
	opts := urlfetch.DefaultConversionOptions()
	fetcher := urlfetch.NewWebFetcherWithConfig(cfg)
	finalURL, html, fetchErr := fetcher.FetchWithURL(url)
	if fetchErr != nil {
		return nil, editor.NewBridgeError(editor.ErrCodeOpFailed, fetchErr.Error())
	}

	parser := urlfetch.NewWebParserWithConfig(cfg)
	webContent, parseErr := parser.Parse(html, finalURL)
	if parseErr != nil {
		return nil, editor.NewBridgeError(editor.ErrCodeOpFailed, parseErr.Error())
	}

	converter := urlfetch.NewURLFetchConverterWithConfig(cfg)
	slides, convErr := converter.ConvertToSlides(webContent, &opts)
	if convErr != nil {
		return nil, editor.NewBridgeError(editor.ErrCodeOpFailed, convErr.Error())
	}

	firstIndex := -1
	for i, slide := range slides {
		idx, addErr := e.AddSlide(slide)
		if addErr != nil {
			return nil, editor.NewBridgeError(editor.ErrCodeOpFailed, addErr.Error())
		}
		if i == 0 {
			firstIndex = idx
		}
	}

	return map[string]int{
		"slide_count": len(slides),
		"first_index": firstIndex,
	}, nil
}

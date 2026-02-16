package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestFindAndReplaceInShapes(t *testing.T) {
	path := writeDeckFixture(t, "find-replace.pptx", []elements.SlideContent{
		elements.NewSlide("Hello World").AddBullet("Hello from bullets"),
		elements.NewSlide("Another Hello").AddBullet("No-op"),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	count, err := editor.FindAndReplaceInShapes("Hello", "Hi")
	if err != nil {
		t.Fatalf("find/replace failed: %v", err)
	}
	if count < 2 {
		t.Fatalf("expected at least 2 replacements, got %d", count)
	}

	slide1, _ := editor.parts.Get("ppt/slides/slide1.xml")
	if !strings.Contains(string(slide1), "Hi World") {
		t.Fatalf("expected replaced title text in slide1")
	}
}

func TestFindAndReplaceInShapesHandlesTextRunsWithAttributes(t *testing.T) {
	path := writeDeckFixture(t, "find-replace-attrs.pptx", []elements.SlideContent{
		elements.NewSlide("Hello World"),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	slide1, ok := editor.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatalf("expected slide part")
	}
	updatedSlide := strings.Replace(
		string(slide1),
		"<a:t>Hello World</a:t>",
		`<a:t xml:space="preserve">Hello World</a:t>`,
		1,
	)
	editor.parts.Set("ppt/slides/slide1.xml", []byte(updatedSlide))

	count, err := editor.FindAndReplaceInShapes("Hello", "Hi")
	if err != nil {
		t.Fatalf("find/replace failed: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected at least one replacement, got %d", count)
	}
	slide1, _ = editor.parts.Get("ppt/slides/slide1.xml")
	if !strings.Contains(string(slide1), `<a:t xml:space="preserve">Hi World</a:t>`) {
		t.Fatalf("expected replacement to preserve a:t attributes")
	}
}

func TestSwapImageByIndex(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, testutil.TinyPNG, 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	path := writeDeckFixture(t, "swap-image.pptx", []elements.SlideContent{
		elements.NewSlide("Image").AddImage(shapes.NewImage(imgPath, 914400, 914400, 1828800, 1828800)),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	before, err := editor.ListSlideImages(0)
	if err != nil {
		t.Fatalf("list images before: %v", err)
	}
	if len(before) != 1 {
		t.Fatalf("expected one image ref, got %d", len(before))
	}
	oldTarget := before[0].Target

	if swapErr := editor.SwapImageByIndex(0, 0, []byte("replacement-image-bytes"), "png"); swapErr != nil {
		t.Fatalf("swap image failed: %v", swapErr)
	}
	after, err := editor.ListSlideImages(0)
	if err != nil {
		t.Fatalf("list images after: %v", err)
	}
	if after[0].Target == oldTarget {
		t.Fatalf("expected image target to change after swap")
	}
}

func TestSwapImageByRelIDDoesNotRegisterImageOnMissingRelationship(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, testutil.TinyPNG, 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	path := writeDeckFixture(t, "swap-image-missing-rel.pptx", []elements.SlideContent{
		elements.NewSlide("Image").AddImage(shapes.NewImage(imgPath, 914400, 914400, 1828800, 1828800)),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	before := len(editor.parts.KeysWithPrefix("ppt/media/image"))
	err = editor.SwapImageByRelID(0, "rId999", []byte("replacement-image-bytes"), "png")
	if err == nil {
		t.Fatalf("expected missing image relationship error")
	}
	after := len(editor.parts.KeysWithPrefix("ppt/media/image"))
	if after != before {
		t.Fatalf("expected media part count unchanged on failure, before=%d after=%d", before, after)
	}
}

func TestSearchShapes(t *testing.T) {
	path := writeDeckFixture(t, "shape-search.pptx", []elements.SlideContent{
		elements.NewSlide("Search Deck"),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	shapeID, err := editor.AddShape(0, "rect", 100, 100, 1000, 600)
	if err != nil {
		t.Fatalf("add shape: %v", err)
	}
	text := "Needle Text"
	if updateErr := editor.UpdateShape(0, shapeID, common.ShapeUpdate{Text: &text}); updateErr != nil {
		t.Fatalf("update shape: %v", updateErr)
	}

	results, err := editor.SearchShapes(common.ShapeSearchQuery{TextContains: "needle", CaseSensitive: false})
	if err != nil {
		t.Fatalf("search shapes: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected at least one matched shape")
	}
}

func TestMergeFromFileWithImageAndChart(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "src.png")
	if err := os.WriteFile(imgPath, testutil.TinyPNG, 0o600); err != nil {
		t.Fatalf("write source image fixture: %v", err)
	}

	sourcePath := writeDeckFixture(t, "source-mixed.pptx", []elements.SlideContent{
		elements.NewSlide("Source Mixed").AddImage(shapes.NewImage(imgPath, 0, 0, 914400, 914400)),
	})
	sourceEditor, err := OpenPresentationEditor(sourcePath)
	if err != nil {
		t.Fatalf("open source editor: %v", err)
	}
	chartDef := charts.NewBarChart([]string{"A", "B"}, []float64{1, 2}).WithTitle("C")
	if addChartErr := sourceEditor.AddChart(0, chartDef); addChartErr != nil {
		t.Fatalf("add chart to source: %v", addChartErr)
	}
	sourceOut := filepath.Join(tmpDir, "source-mixed-chart.pptx")
	if saveSourceErr := sourceEditor.Save(sourceOut); saveSourceErr != nil {
		t.Fatalf("save source with chart: %v", saveSourceErr)
	}
	_ = sourceEditor.Close()

	destPath := writeDeckFixture(t, "dest-mixed.pptx", []elements.SlideContent{
		elements.NewSlide("Dest"),
	})
	destEditor, err := OpenPresentationEditor(destPath)
	if err != nil {
		t.Fatalf("open dest editor: %v", err)
	}
	defer func() { _ = destEditor.Close() }()

	if mergeErr := destEditor.MergeFromFile(sourceOut); mergeErr != nil {
		t.Fatalf("merge from mixed source failed: %v", mergeErr)
	}

	outPath := filepath.Join(tmpDir, "merged-mixed.pptx")
	if saveDestErr := destEditor.Save(outPath); saveDestErr != nil {
		t.Fatalf("save merged mixed deck: %v", saveDestErr)
	}

	reopened, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen merged deck: %v", err)
	}
	defer func() { _ = reopened.Close() }()
	if reopened.SlideCount() != 2 {
		t.Fatalf("expected 2 slides, got %d", reopened.SlideCount())
	}
	if len(reopened.parts.KeysWithPrefix("ppt/media/image")) == 0 {
		t.Fatalf("expected merged image media parts")
	}
	if len(reopened.parts.KeysWithPrefix("ppt/charts/chart")) == 0 {
		t.Fatalf("expected merged chart parts")
	}
	if len(reopened.parts.KeysWithPrefix("ppt/embeddings/")) == 0 {
		t.Fatalf("expected merged embedding parts")
	}
}

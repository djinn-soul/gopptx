// examples/79-document-infra demonstrates document infrastructure operations:
// sections, slide comments, z-order (bring to front / send to back), and
// shape grouping/ungrouping via the PresentationEditor.
//
// Run with: go run ./examples/79-document-infra/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	outputDir  = "examples/output"
	outputFile = "79_document_infra.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-infra-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	basePath := filepath.Join(tmpDir, "base.pptx")
	if err := buildBasePresentation(basePath); err != nil {
		return err
	}

	ed, err := editor.OpenPresentationEditor(basePath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer ed.Close()

	if err := addAndRenameSections(ed); err != nil {
		return err
	}

	outPath := filepath.Join(outputDir, outputFile)
	applyZOrderAndGrouping(ed, 4)
	if err := ed.Save(outPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return verifyOutput(outPath)
}

func buildBasePresentation(basePath string) error {
	rect1 := pptx.NewRectangle(1, 1.5, 3, 1.5).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Bottom shape (z-order)").
		WithName("rect-bottom")

	rect2 := pptx.NewRectangle(2, 2, 3, 1.5).
		WithFill(pptx.NewShapeFill("C0504D")).
		WithText("Top shape (z-order)").
		WithName("rect-top")

	builder := pptx.NewPresentationBuilder("Document Infrastructure Demo")
	builder.AddSlide(pptx.NewSlide("Section A – Slide 1").AddBullet("First slide in Section A"))
	builder.AddSlide(pptx.NewSlide("Section A – Slide 2").AddBullet("Second slide in Section A"))
	builder.AddSlide(pptx.NewSlide("Section B – Slide 3").AddBullet("First slide in Section B"))
	builder.AddSlide(pptx.NewSlide("Section B – Slide 4").AddBullet("Second slide in Section B"))

	appendixSlide := pptx.NewSlide("Appendix – Z-Order Demo").
		AddShape(rect1).
		AddShape(rect2).
		AddBullet("rect-bottom is added first; rect-top overlaps it.")
	builder.AddSlide(appendixSlide)

	commentSlide := pptx.NewSlide("Slide with Comments").
		AddBullet("This slide has author comments.").
		AddComment("Alice", "Great point – add more data here.").
		AddComment("Bob", "Needs a chart to illustrate this.")
	builder.AddSlide(commentSlide)

	if err := builder.WriteToFile(basePath); err != nil {
		return fmt.Errorf("write base: %w", err)
	}
	return nil
}

func addAndRenameSections(ed *editor.PresentationEditor) error {
	log.Println("Adding sections...")
	if err := ed.AddSection("Introduction", []int{0, 1}); err != nil {
		return fmt.Errorf("add introduction section: %w", err)
	}
	if err := ed.AddSection("Core Content", []int{2, 3}); err != nil {
		return fmt.Errorf("add core content section: %w", err)
	}
	if err := ed.AddSection("Appendix", []int{4, 5}); err != nil {
		return fmt.Errorf("add appendix section: %w", err)
	}

	log.Println("Renaming 'Appendix' to 'Back Matter'...")
	if err := ed.RenameSection("Appendix", "Back Matter"); err != nil {
		return fmt.Errorf("rename section: %w", err)
	}

	logSections(ed)
	return nil
}

func logSections(ed *editor.PresentationEditor) {
	sections := ed.Sections()
	log.Printf("Sections after add/rename: %d\n", len(sections))
	for _, s := range sections {
		log.Printf("  - %q (slides: %d)\n", s.Name, len(s.SlideIDs))
	}
}

func applyZOrderAndGrouping(ed *editor.PresentationEditor, slideIdx int) {
	shapeList, err := ed.GetShapes(slideIdx)
	if err != nil {
		log.Printf("get shapes on slide %d: %v (non-fatal)\n", slideIdx, err)
		return
	}
	log.Printf("Shapes on slide %d: %d\n", slideIdx, len(shapeList))

	applyZOrder(ed, slideIdx, shapeList)
	applyGrouping(ed, slideIdx, shapeList)
}

func applyZOrder(ed *editor.PresentationEditor, slideIdx int, shapeList []editorcommon.Shape) {
	if len(shapeList) < 2 {
		return
	}
	firstID := shapeList[0].ID
	log.Printf("Bringing shape ID %d to front...\n", firstID)
	if err := ed.BringShapeToFront(slideIdx, firstID); err != nil {
		log.Printf("BringShapeToFront: %v (non-fatal)\n", err)
	}

	lastID := shapeList[len(shapeList)-1].ID
	log.Printf("Sending shape ID %d to back...\n", lastID)
	if err := ed.SendShapeToBack(slideIdx, lastID); err != nil {
		log.Printf("SendShapeToBack: %v (non-fatal)\n", err)
	}
}

func applyGrouping(ed *editor.PresentationEditor, slideIdx int, shapeList []editorcommon.Shape) {
	if len(shapeList) < 2 {
		return
	}
	ids := make([]int, len(shapeList))
	for i, s := range shapeList {
		ids[i] = s.ID
	}
	log.Printf("Grouping %d shapes on slide %d...\n", len(ids), slideIdx)
	groupID, err := ed.GroupShapes(slideIdx, ids)
	if err != nil {
		log.Printf("GroupShapes: %v (non-fatal)\n", err)
		return
	}
	log.Printf("Group ID: %d. Ungrouping...\n", groupID)
	if _, err := ed.UngroupShapes(slideIdx, groupID); err != nil {
		log.Printf("UngroupShapes: %v (non-fatal)\n", err)
	}
}

func verifyOutput(outPath string) error {
	ed2, err := editor.OpenPresentationEditor(outPath)
	if err != nil {
		return fmt.Errorf("reopen for verification: %w", err)
	}
	defer ed2.Close()

	secs := ed2.Sections()
	log.Printf("Verified sections in output: %d\n", len(secs))
	for _, s := range secs {
		log.Printf("  - %q\n", s.Name)
	}

	log.Printf("Generated %s\n", outPath)
	return nil
}

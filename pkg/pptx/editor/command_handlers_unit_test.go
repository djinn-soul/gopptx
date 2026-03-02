package editor

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCommandHandlers_Content(t *testing.T) {
	fixturePath := filepath.Join(testutil.RootTestdataDir(), "simple.pptx")
	e, err := OpenPresentationEditor(fixturePath)
	if err != nil {
		t.Fatalf("failed to open editor: %v", err)
	}
	defer e.Close()

	t.Run("Authors", func(t *testing.T) {
		payload := []byte(`{"name": "Test Author", "initials": "TA"}`)
		res, err := handleAddAuthor(e, payload)
		if err != nil {
			t.Fatalf("handleAddAuthor failed: %v", err)
		}

		authorsResp, err := handleGetAuthors(e, nil)
		if err != nil {
			t.Fatalf("handleGetAuthors failed: %v", err)
		}

		resMap := res.(map[string]int64)
		if resMap["author_id"] == 0 {
			t.Error("expected valid author_id")
		}

		if authorsResp == nil {
			t.Error("expected authors list")
		}
	})

	t.Run("CustomXML", func(t *testing.T) {
		payload := []byte(`{"content": "<test/>", "properties": {"key": "val"}}`)
		res, err := handleAddCustomXML(e, payload)
		if err != nil {
			t.Fatalf("handleAddCustomXML failed: %v", err)
		}

		listResp, err := handleListCustomXML(e, nil)
		if err != nil {
			t.Fatalf("handleListCustomXML failed: %v", err)
		}
		if listResp == nil { t.Error("expected custom xml list") }

		resMap := res.(map[string]int)
		idx := resMap["index"]

		remPayload := []byte(`{"index": ` + strconv.Itoa(idx) + `}`)
		_, err = handleRemoveCustomXML(e, remPayload)
		if err != nil {
			t.Fatalf("handleRemoveCustomXML failed: %v", err)
		}
	})

	t.Run("Protection", func(t *testing.T) {
		p1 := []byte(`{"password": "secret"}`)
		_, err := handleSetModifyPassword(e, p1)
		if err != nil { t.Error(err) }

		p2 := []byte(`{"final": true}`)
		_, err = handleSetMarkAsFinal(e, p2)
		if err != nil { t.Error(err) }
	})

	t.Run("Slides", func(t *testing.T) {
		payload := []byte(`{"title": "New Slide", "layout": "titleAndContent"}`)
		_, err := handleAddSlide(e, payload)
		if err != nil {
			t.Fatalf("handleAddSlide failed: %v", err)
		}

		listResp, err := handleListSlides(e, nil)
		if err != nil { t.Fatalf("handleListSlides failed: %v", err) }
		if listResp == nil { t.Error("expected slide list") }
	})

	t.Run("Shapes", func(t *testing.T) {
		// Use slide index 1 (the one we just created)
		payload := []byte(`{"slide_index": 1, "type": "rect", "x": 100.0, "y": 100.0, "w": 100.0, "h": 100.0, "text": "Hello"}`)
		res, err := handleAddShape(e, payload)
		if err != nil {
			t.Fatalf("handleAddShape failed: %v", err)
		}

		resMap := res.(map[string]int)
		shapeID := resMap["shape_id"]
		if shapeID == 0 { t.Error("expected valid shape_id") }

		listPayload := []byte(`{"slide_index": 1}`)
		listResp, err := handleListShapes(e, listPayload)
		if err != nil { t.Fatalf("handleListShapes failed: %v", err) }
		if listResp == nil { t.Error("expected shape list") }

		// Update Shape
		upPayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(shapeID) + `, "updates": {"text": "New Text"}}`)
		_, err = handleUpdateShape(e, upPayload)
		if err != nil { t.Fatalf("handleUpdateShape failed: %v", err) }

		// Search Shapes
		searchPayload := []byte(`{"text_contains": "New Text"}`)
		_, err = handleSearchShapes(e, searchPayload)
		if err != nil { t.Fatalf("handleSearchShapes failed: %v", err) }

		// Move Shapes
		movePayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(shapeID) + `}`)
		_, _ = handleMoveShapeToFront(e, movePayload)
		_, _ = handleMoveShapeToBack(e, movePayload)

		// Remove Shape
		_, err = handleRemoveShape(e, movePayload)
		if err != nil { t.Fatalf("handleRemoveShape failed: %v", err) }

		// Add Image
		// Use a dummy file for the path check
		tmpDir := t.TempDir()
		dummyImg := filepath.Join(tmpDir, "dummy.png")
		_ = os.WriteFile(dummyImg, []byte("fake-png-data"), 0600)

		imgPayload := []byte(`{"slide_index": 1, "path": "` + filepath.ToSlash(dummyImg) + `", "x": 0.0, "y": 0.0, "w": 50.0, "h": 50.0}`)
		_, _ = handleAddImage(e, imgPayload) // ignore error as fake-png-data might fail decoding, but it hits the path

		imgBytesNoFormatPayload := []byte(`{"slide_index": 1, "data": "AQID", "x": 0.0, "y": 0.0, "w": 50.0, "h": 50.0}`)
		_, err = handleAddImage(e, imgBytesNoFormatPayload)
		if err == nil {
			t.Fatal("expected handleAddImage to reject byte payload without format")
		}
	})

	t.Run("CommentsAndNotes", func(t *testing.T) {
		// Needs an author first
		authPayload := []byte(`{"name": "Commenter", "initials": "C"}`)
		_, _ = handleAddAuthor(e, authPayload)

		// Add comment
		addComPayload := []byte(`{"slide_index": 1, "author_id": 2, "text": "Comment", "x": 10, "y": 10}`)
		_, _ = handleAddComment(e, addComPayload)

		// Get comments
		getComPayload := []byte(`{"slide_index": 1}`)
		_, _ = handleGetComments(e, getComPayload)

		// Notes
		setNotesPayload := []byte(`{"slide_index": 1, "text": "Notes text"}`)
		_, _ = handleSetNotes(e, setNotesPayload)

		getNotesPayload := []byte(`{"slide_index": 1}`)
		_, _ = handleGetNotes(e, getNotesPayload)
	})

	t.Run("Tables", func(t *testing.T) {
		// Add Table
		payload := []byte(`{"slide_index": 1, "x": 100, "y": 100, "rows": 2, "cols": 2}`)
		res, err := handleAddTable(e, payload)
		if err != nil { t.Fatalf("handleAddTable failed: %v", err) }

		resMap := res.(map[string]int)
		tableID := resMap["shape_id"]

		// Get Table
		getPayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(tableID) + `}`)
		_, err = handleGetTable(e, getPayload)
		if err != nil { t.Fatalf("handleGetTable failed: %v", err) }

		// Update Table Cell
		updPayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(tableID) + `, "row": 0, "col": 0, "updates": {"text": "Cell 0,0"}}`)
		_, err = handleUpdateTableCell(e, updPayload)
		if err != nil { t.Fatalf("handleUpdateTableCell failed: %v", err) }

		// Update Table Flags
		flagsPayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(tableID) + `, "flags": {"first_row": true}}`)
		_, err = handleUpdateTableFlags(e, flagsPayload)
		if err != nil { t.Fatalf("handleUpdateTableFlags failed: %v", err) }

		// Merge Cells
		mergePayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(tableID) + `, "row1": 0, "col1": 0, "row2": 1, "col2": 1}`)
		_, err = handleMergeTableCells(e, mergePayload)
		if err != nil { t.Fatalf("handleMergeTableCells failed: %v", err) }

		// Split Cell
		splitPayload := []byte(`{"slide_index": 1, "shape_id": ` + strconv.Itoa(tableID) + `, "row": 0, "col": 0}`)
		_, err = handleSplitTableCell(e, splitPayload)
		if err != nil { t.Fatalf("handleSplitTableCell failed: %v", err) }
	})
}

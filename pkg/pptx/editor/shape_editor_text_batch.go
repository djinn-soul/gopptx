package editor

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// UpdateSlideRunTexts updates multiple run texts on a slide in one parse/rewrite pass.
func (e *PresentationEditor) UpdateSlideRunTexts(
	slideIndex int,
	updates []common.ShapeRunTextUpdate,
) error {
	if len(updates) == 0 {
		return nil
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	updatesByShape := make(map[int][]common.ShapeRunTextUpdate, len(updates))
	for _, update := range updates {
		updatesByShape[update.ShapeID] = append(updatesByShape[update.ShapeID], update)
	}

	found := make(map[int]bool, len(updatesByShape))
	var updateErr error
	newXML := replaceShapeNodes(content, shapes, func(_ int, shape *parsedShape) ([]byte, bool) {
		targetID, ok := matchRunTextTarget(shape, updatesByShape)
		if !ok || updateErr != nil {
			return nil, false
		}

		runs := editorshape.CopyTextRuns(shape.Runs)
		for _, update := range updatesByShape[targetID] {
			runs, updateErr = editorshape.UpdateRunText(runs, update.RunIndex, update.Text)
			if updateErr != nil {
				return nil, false
			}
		}
		shape.Runs = runs

		updatedXML, err := replaceShapeTextBody(e, partPath, content[shape.Start:shape.End], shape)
		if err != nil {
			updateErr = err
			return nil, false
		}
		found[targetID] = true
		return updatedXML, true
	})
	if updateErr != nil {
		return updateErr
	}
	for shapeID := range updatesByShape {
		if !found[shapeID] {
			return fmt.Errorf("shape id %d not found on slide %d", shapeID, slideIndex)
		}
	}

	e.parts.Set(partPath, newXML)
	return nil
}

// UpdateDeckRunTexts updates run texts across multiple slides.
func (e *PresentationEditor) UpdateDeckRunTexts(
	slideUpdates []common.SlideRunTextUpdates,
) error {
	if len(slideUpdates) == 0 {
		return nil
	}

	groupedBySlide := make(map[int][]common.ShapeRunTextUpdate, len(slideUpdates))
	slideOrder := make([]int, 0, len(slideUpdates))
	for _, slideUpdate := range slideUpdates {
		if _, seen := groupedBySlide[slideUpdate.SlideIndex]; !seen {
			slideOrder = append(slideOrder, slideUpdate.SlideIndex)
		}
		groupedBySlide[slideUpdate.SlideIndex] = append(
			groupedBySlide[slideUpdate.SlideIndex],
			slideUpdate.Updates...,
		)
	}
	if len(slideOrder) == 1 {
		slideIndex := slideOrder[0]
		return e.UpdateSlideRunTexts(slideIndex, groupedBySlide[slideIndex])
	}

	workerCount := runtime.GOMAXPROCS(0)
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > len(slideOrder) {
		workerCount = len(slideOrder)
	}
	type job struct {
		slideIndex int
	}
	jobs := make(chan job, len(slideOrder))
	errBySlide := make(map[int]error, len(slideOrder))
	var (
		errMu sync.Mutex
		wg    sync.WaitGroup
	)
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				if err := e.UpdateSlideRunTexts(item.slideIndex, groupedBySlide[item.slideIndex]); err != nil {
					errMu.Lock()
					errBySlide[item.slideIndex] = err
					errMu.Unlock()
				}
			}
		}()
	}
	for _, slideIndex := range slideOrder {
		jobs <- job{slideIndex: slideIndex}
	}
	close(jobs)
	wg.Wait()

	for _, slideIndex := range slideOrder {
		if err := errBySlide[slideIndex]; err != nil {
			return err
		}
	}
	return nil
}

func matchRunTextTarget(
	shape *parsedShape,
	updatesByShape map[int][]common.ShapeRunTextUpdate,
) (int, bool) {
	if _, ok := updatesByShape[shape.ID]; ok {
		return shape.ID, true
	}
	if shape.PhType == placeholderTypeTitle {
		if _, ok := updatesByShape[0]; ok {
			return 0, true
		}
	}
	return 0, false
}

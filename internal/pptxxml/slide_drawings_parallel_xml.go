package pptxxml

import (
	"runtime"
	"sync"
)

const minShapesForParallel = 4

func renderCustomShapeXMLConcurrently(shapes []ShapeSpec, firstShapeID int) ([]int, []string) {
	if len(shapes) == 0 {
		return nil, nil
	}

	shapeIDs := make([]int, len(shapes))
	xmlByIndex := make([]string, len(shapes))
	if len(shapes) <= minShapesForParallel {
		for i := range shapes {
			id := firstShapeID + i
			shapeIDs[i] = id
			xmlByIndex[i] = customShapeXML(shapes[i], id)
		}
		return shapeIDs, xmlByIndex
	}

	workerCount := max(runtime.GOMAXPROCS(0), 1)
	if workerCount > len(shapes) {
		workerCount = len(shapes)
	}

	jobs := make(chan int, workerCount*2)
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for range workerCount {
		go func() {
			defer wg.Done()
			for index := range jobs {
				id := firstShapeID + index
				shapeIDs[index] = id
				xmlByIndex[index] = customShapeXML(shapes[index], id)
			}
		}()
	}

	for i := range shapes {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return shapeIDs, xmlByIndex
}

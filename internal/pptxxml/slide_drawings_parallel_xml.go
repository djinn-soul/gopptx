package pptxxml

import "sync"

func renderCustomShapeXMLConcurrently(shapes []ShapeSpec, firstShapeID int) ([]int, []string) {
	if len(shapes) == 0 {
		return nil, nil
	}

	shapeIDs := make([]int, len(shapes))
	xmlByIndex := make([]string, len(shapes))
	if len(shapes) == 1 {
		shapeIDs[0] = firstShapeID
		xmlByIndex[0] = customShapeXML(shapes[0], firstShapeID)
		return shapeIDs, xmlByIndex
	}

	var wg sync.WaitGroup
	wg.Add(len(shapes))
	for i := range shapes {
		index := i
		go func() {
			defer wg.Done()
			id := firstShapeID + index
			shapeIDs[index] = id
			xmlByIndex[index] = customShapeXML(shapes[index], id)
		}()
	}
	wg.Wait()
	return shapeIDs, xmlByIndex
}

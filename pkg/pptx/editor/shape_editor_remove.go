package editor

func (e *PresentationEditor) applyShapeRemoval(
	partPath string,
	content []byte,
	shapes []parsedShape,
	shapeIndex int,
) error {
	newContent := replaceShapeNodes(content, shapes, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == shapeIndex {
			return []byte{}, true
		}
		return nil, false
	})
	e.parts.Set(partPath, newContent)
	return nil
}

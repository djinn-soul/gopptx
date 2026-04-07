package editor

type updatedResponse struct {
	Updated bool `json:"updated"`
}

type addedResponse struct {
	Added bool `json:"added"`
}

type removedResponse struct {
	Removed bool `json:"removed"`
}

type mergedResponse struct {
	Merged bool `json:"merged"`
}

type appliedResponse struct {
	Applied bool `json:"applied"`
}

type swappedResponse struct {
	Swapped bool `json:"swapped"`
}

type successResponse struct {
	Success bool `json:"success"`
}

type clearedResponse struct {
	Cleared bool `json:"cleared"`
}

type deletedResponse struct {
	Deleted bool `json:"deleted"`
}

type shapeIDResponse struct {
	ShapeID int `json:"shape_id"`
}

type groupIDResponse struct {
	GroupID int `json:"group_id"`
}

type countResponse struct {
	Count int `json:"count"`
}

type indexResponse struct {
	Index int `json:"index"`
}

type newIndexResponse struct {
	NewIndex int `json:"new_index"`
}

type replacementsResponse struct {
	Replacements int `json:"replacements"`
}

type markdownSlidesResponse struct {
	SlideCount int `json:"slide_count"`
	FirstIndex int `json:"first_index"`
}

type mermaidAddResponse struct {
	ShapeCount     int `json:"shape_count"`
	ConnectorCount int `json:"connector_count"`
}

var (
	respUpdated = updatedResponse{Updated: true}
	respAdded   = addedResponse{Added: true}
	respRemoved = removedResponse{Removed: true}
	respMerged  = mergedResponse{Merged: true}
	respApplied = appliedResponse{Applied: true}
	respSwapped = swappedResponse{Swapped: true}
	respSuccess = successResponse{Success: true}
	respCleared = clearedResponse{Cleared: true}
	respDeleted = deletedResponse{Deleted: true}
)

func respShapeID(id int) shapeIDResponse { return shapeIDResponse{ShapeID: id} }
func respGroupID(id int) groupIDResponse { return groupIDResponse{GroupID: id} }
func respCount(v int) countResponse      { return countResponse{Count: v} }
func respIndex(v int) indexResponse      { return indexResponse{Index: v} }
func respNewIndex(v int) newIndexResponse {
	return newIndexResponse{NewIndex: v}
}
func respReplacements(v int) replacementsResponse {
	return replacementsResponse{Replacements: v}
}

package command

import "strconv"

type TextboxPlacementRequest struct {
	SlideIndex int
	Left       float64
	Top        float64
	Width      float64
	Height     float64
}

type GroupShapeRequest struct {
	SlideIndex int
	ShapeIDs   []int
}

func ParseGroupShapeRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
	parseIntSliceField ParseIntSliceFieldFn,
) (GroupShapeRequest, bool) {
	slideIndex, ok := parseIntField(payload, "slide_index")
	if !ok {
		return GroupShapeRequest{}, false
	}
	shapeIDs := []int{}
	if _, present := payload["shapes"]; present {
		parsedIDs, idsOK := parseIntSliceField(payload, "shapes")
		if !idsOK {
			return GroupShapeRequest{}, false
		}
		shapeIDs = parsedIDs
	}
	return GroupShapeRequest{
		SlideIndex: slideIndex,
		ShapeIDs:   shapeIDs,
	}, true
}

type FreeformPoint struct {
	X int
	Y int
}

const minFreeformPointCount = 2

func ParseFreeformPoints(payload map[string]any) ([]FreeformPoint, error) {
	rawPoints, exists := payload["points"]
	if !exists {
		return nil, NewValidationError("MISSING_FIELD", "missing required field: points")
	}
	pointItems, ok := rawPoints.([]any)
	if !ok {
		return nil, NewValidationError("INVALID_TYPE", "field points must be an array")
	}
	if len(pointItems) < minFreeformPointCount {
		return nil, NewValidationError("INVALID_VALUE", "field points must contain at least 2 points")
	}

	points := make([]FreeformPoint, 0, len(pointItems))
	for i, item := range pointItems {
		pair, pairOK := item.([]any)
		if !pairOK || len(pair) != 2 {
			return nil, NewValidationError("INVALID_TYPE", "field points["+strconv.Itoa(i)+"] must be [x, y]")
		}
		x, xOK := numberToFloat64(pair[0])
		y, yOK := numberToFloat64(pair[1])
		if !xOK || !yOK {
			return nil, NewValidationError(
				"INVALID_TYPE",
				"field points["+strconv.Itoa(i)+"] coordinates must be numbers",
			)
		}
		points = append(points, FreeformPoint{X: int(x), Y: int(y)})
	}
	return points, nil
}

func ParseOptionalCloseFlag(payload map[string]any) (bool, error) {
	rawClose, exists := payload["close"]
	if !exists {
		return false, nil
	}
	flag, castOK := rawClose.(bool)
	if !castOK {
		return false, NewValidationError("INVALID_TYPE", "field close must be a boolean")
	}
	return flag, nil
}

func numberToFloat64(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	default:
		return 0, false
	}
}

func ParseTextboxPlacementRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseFloatField RequireFloatFn,
) (TextboxPlacementRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return TextboxPlacementRequest{}, false
	}
	left, ok := parseFloatField(payload, "left")
	if !ok {
		return TextboxPlacementRequest{}, false
	}
	top, ok := parseFloatField(payload, "top")
	if !ok {
		return TextboxPlacementRequest{}, false
	}
	width, ok := parseFloatField(payload, "width")
	if !ok {
		return TextboxPlacementRequest{}, false
	}
	height, ok := parseFloatField(payload, "height")
	if !ok {
		return TextboxPlacementRequest{}, false
	}
	return TextboxPlacementRequest{
		SlideIndex: slideIndex,
		Left:       left,
		Top:        top,
		Width:      width,
		Height:     height,
	}, true
}

type ConnectorPlacementRequest struct {
	SlideIndex    int
	ConnectorType string
	BeginX        float64
	BeginY        float64
	EndX          float64
	EndY          float64
}

func ParseConnectorPlacementRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseStringField ParseStringFieldFn,
	parseFloatField RequireFloatFn,
) (ConnectorPlacementRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return ConnectorPlacementRequest{}, false
	}
	connectorType, ok := parseStringField(payload, "connector_type")
	if !ok {
		return ConnectorPlacementRequest{}, false
	}
	beginX, ok := parseFloatField(payload, "begin_x")
	if !ok {
		return ConnectorPlacementRequest{}, false
	}
	beginY, ok := parseFloatField(payload, "begin_y")
	if !ok {
		return ConnectorPlacementRequest{}, false
	}
	endX, ok := parseFloatField(payload, "end_x")
	if !ok {
		return ConnectorPlacementRequest{}, false
	}
	endY, ok := parseFloatField(payload, "end_y")
	if !ok {
		return ConnectorPlacementRequest{}, false
	}
	return ConnectorPlacementRequest{
		SlideIndex:    slideIndex,
		ConnectorType: connectorType,
		BeginX:        beginX,
		BeginY:        beginY,
		EndX:          endX,
		EndY:          endY,
	}, true
}

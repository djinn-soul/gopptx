package command

type OptionalStringSliceFieldFn func(map[string]any, string) ([]string, bool)
type OptionalIntFieldFn func(map[string]any, string) (int, bool)

type AddSlideRequest struct {
	Title   string
	Layout  string
	Bullets []string
}

func ParseAddSlideRequest(
	payload map[string]any,
	optionalString OptionalStringFieldFn,
	optionalStringSlice OptionalStringSliceFieldFn,
) AddSlideRequest {
	bullets, _ := optionalStringSlice(payload, "bullets")
	return AddSlideRequest{
		Title:   optionalString(payload, "title"),
		Layout:  optionalString(payload, "layout"),
		Bullets: bullets,
	}
}

type SlideIndexRequest struct {
	Index int
}

func ParseSlideIndexOnlyRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
) (SlideIndexRequest, bool) {
	index, ok := parseIntField(payload, "index")
	if !ok {
		return SlideIndexRequest{}, false
	}
	return SlideIndexRequest{Index: index}, true
}

type MoveSlideRequest struct {
	From int
	To   int
}

func ParseMoveSlideRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
) (MoveSlideRequest, bool) {
	from, ok := parseIntField(payload, "from")
	if !ok {
		return MoveSlideRequest{}, false
	}
	to, ok := parseIntField(payload, "to")
	if !ok {
		return MoveSlideRequest{}, false
	}
	return MoveSlideRequest{
		From: from,
		To:   to,
	}, true
}

type DuplicateSlideRequest struct {
	Index    int
	InsertAt int
}

func ParseDuplicateSlideRequest(
	payload map[string]any,
	parseIntField ParseIntFieldFn,
	optionalIntField OptionalIntFieldFn,
) (DuplicateSlideRequest, bool) {
	index, ok := parseIntField(payload, "index")
	if !ok {
		return DuplicateSlideRequest{}, false
	}
	insertAt, _ := optionalIntField(payload, "insert_at")
	return DuplicateSlideRequest{
		Index:    index,
		InsertAt: insertAt,
	}, true
}

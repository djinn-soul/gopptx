package command

type MediaPlacement struct {
	SlideIndex int
	X          float64
	Y          float64
	W          float64
	H          float64
}

type MediaInsertSpec struct {
	MetaKey          string
	PrimaryPathKey   string
	PrimaryDataKey   string
	SecondaryPathKey string
	SecondaryDataKey string
	PrimaryMaxLen    int
	SecondaryMaxLen  int
	PrimaryLabel     string
	SecondaryLabel   string
	InsertBinary     func(MediaPlacement, string, []byte, []byte) (int, error)
	InsertPath       func(MediaPlacement, string, string, string) (int, error)
}

type (
	RequireIntFn      func(map[string]any, string) (int, bool)
	RequireFloatFn    func(map[string]any, string) (float64, bool)
	IndexBoundsFn     func(int, int, int, string) bool
	ParseRawPayloadFn func([]byte) (map[string]any, error)
	BuildResultFn     func(int) any
)

func ExecuteMediaInsert(
	payload map[string]any,
	placement MediaPlacement,
	optionalString func(map[string]any, string) string,
	spec MediaInsertSpec,
) (int, error) {
	metaValue := optionalString(payload, spec.MetaKey)
	primaryPath := optionalString(payload, spec.PrimaryPathKey)
	primaryBase64 := optionalString(payload, spec.PrimaryDataKey)
	secondaryPath := optionalString(payload, spec.SecondaryPathKey)
	secondaryBase64 := optionalString(payload, spec.SecondaryDataKey)

	primaryData, err := DecodeOptionalBase64Field(primaryBase64, spec.PrimaryMaxLen, spec.PrimaryLabel)
	if err != nil {
		return 0, err
	}
	secondaryData, err := DecodeOptionalBase64Field(secondaryBase64, spec.SecondaryMaxLen, spec.SecondaryLabel)
	if err != nil {
		return 0, err
	}

	return insertShapeFromBinaryOrPath(
		len(primaryData) > 0 || len(secondaryData) > 0,
		func() (int, error) {
			return spec.InsertBinary(placement, metaValue, primaryData, secondaryData)
		},
		func() (int, error) {
			return spec.InsertPath(placement, metaValue, primaryPath, secondaryPath)
		},
	)
}

func HandleMediaInsertCommand(
	rawPayload []byte,
	slideCount int,
	parseRawPayload ParseRawPayloadFn,
	requireInt RequireIntFn,
	requireFloat RequireFloatFn,
	indexBounds IndexBoundsFn,
	optionalString func(map[string]any, string) string,
	validationErr func() error,
	buildResult BuildResultFn,
	spec MediaInsertSpec,
) (any, error) {
	p, err := parseRawPayload(rawPayload)
	if err != nil {
		return nil, err
	}

	placement, ok := ParseMediaPlacement(p, slideCount, requireInt, requireFloat, indexBounds)
	if !ok {
		return nil, validationErr()
	}

	shapeID, err := ExecuteMediaInsert(p, placement, optionalString, spec)
	if err != nil {
		return nil, err
	}
	return buildResult(shapeID), nil
}

func ParseMediaPlacement(
	payload map[string]any,
	slideCount int,
	requireInt RequireIntFn,
	requireFloat RequireFloatFn,
	indexBounds IndexBoundsFn,
) (MediaPlacement, bool) {
	slideIndex, ok := requireInt(payload, "slide_index")
	if !ok {
		return MediaPlacement{}, false
	}
	x, ok := requireFloat(payload, "x")
	if !ok {
		return MediaPlacement{}, false
	}
	y, ok := requireFloat(payload, "y")
	if !ok {
		return MediaPlacement{}, false
	}
	w, ok := requireFloat(payload, "w")
	if !ok {
		return MediaPlacement{}, false
	}
	h, ok := requireFloat(payload, "h")
	if !ok {
		return MediaPlacement{}, false
	}
	if !indexBounds(slideIndex, 0, slideCount, "slide_index") {
		return MediaPlacement{}, false
	}
	return MediaPlacement{
		SlideIndex: slideIndex,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
	}, true
}

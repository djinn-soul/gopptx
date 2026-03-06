package command

import (
	"errors"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type AddImageRequireSlideIndexFn func(map[string]any) (int, bool)
type AddImageRequireFloatFn func(map[string]any, string) (float64, bool)
type AddImageOptionalStringFn func(map[string]any, string) string
type AddImageFromBytesFn func(
	slideIndex int,
	data []byte,
	format string,
	x, y, w, h float64,
	options *common.ShapeUpdate,
) (int, error)
type AddImageFromURLFn func(
	slideIndex int,
	sourceURL string,
	x, y, w, h float64,
	options *common.ShapeUpdate,
) (int, error)
type AddImageFromPathFn func(
	slideIndex int,
	imagePath string,
	x, y, w, h float64,
	options *common.ShapeUpdate,
) (int, error)

type AddImageRequest struct {
	SlideIndex int
	X          float64
	Y          float64
	W          float64
	H          float64
	ImagePath  string
	ImageURL   string
	Base64Data string
	Format     string
	Options    *common.ShapeUpdate
}

func ParseAddImageRequest(
	payload map[string]any,
	requireSlideIndex AddImageRequireSlideIndexFn,
	requireFloat AddImageRequireFloatFn,
	optionalString AddImageOptionalStringFn,
) (AddImageRequest, bool, error) {
	slideIndex, ok := requireSlideIndex(payload)
	if !ok {
		return AddImageRequest{}, false, nil
	}
	x, ok := requireFloat(payload, "x")
	if !ok {
		return AddImageRequest{}, false, nil
	}
	y, ok := requireFloat(payload, "y")
	if !ok {
		return AddImageRequest{}, false, nil
	}
	w, ok := requireFloat(payload, "w")
	if !ok {
		return AddImageRequest{}, false, nil
	}
	h, ok := requireFloat(payload, "h")
	if !ok {
		return AddImageRequest{}, false, nil
	}

	request := AddImageRequest{
		SlideIndex: slideIndex,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
		ImagePath:  optionalString(payload, "path"),
		ImageURL:   optionalString(payload, "url"),
		Base64Data: optionalString(payload, "data"),
		Format:     optionalString(payload, "format"),
	}
	if err := DecodeOptionalPayloadValue(payload, "options", &request.Options); err != nil {
		return AddImageRequest{}, true, err
	}
	return request, true, nil
}

func DecodeImagePayload(base64Data, format string, maxLen int) ([]byte, error) {
	if base64Data == "" {
		return nil, nil
	}
	if format == "" {
		return nil, errors.New("image format is required when image data is provided")
	}
	return DecodeOptionalBase64Field(base64Data, maxLen, "image")
}

func ExecuteAddImageRequest(
	request AddImageRequest,
	maxLen int,
	addImageFromBytes AddImageFromBytesFn,
	addImageFromURL AddImageFromURLFn,
	addImageFromPath AddImageFromPathFn,
) (int, error) {
	if request.Base64Data != "" {
		decodedData, err := DecodeImagePayload(
			request.Base64Data,
			strings.TrimSpace(request.Format),
			maxLen,
		)
		if err != nil {
			return 0, err
		}
		return addImageFromBytes(
			request.SlideIndex,
			decodedData,
			request.Format,
			request.X,
			request.Y,
			request.W,
			request.H,
			request.Options,
		)
	}
	if request.ImageURL != "" {
		return addImageFromURL(
			request.SlideIndex,
			request.ImageURL,
			request.X,
			request.Y,
			request.W,
			request.H,
			request.Options,
		)
	}
	return addImageFromPath(
		request.SlideIndex,
		request.ImagePath,
		request.X,
		request.Y,
		request.W,
		request.H,
		request.Options,
	)
}

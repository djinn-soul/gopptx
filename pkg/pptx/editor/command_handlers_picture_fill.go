package editor

import (
	"encoding/json"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/handlers/slidescatalog"
)

// handleSetPictureFill writes an <a:blipFill> onto a shape or a table cell.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "shape_id": N,
//	  "image_path": "logo.png",     // or "image_data": base64
//	  "mode": "stretch",            // optional: "stretch" (default) | "tile"
//	  "crop": {"left":0.1,"top":0,"right":0.1,"bottom":0}, // optional
//	  "row": 1, "column": 2         // optional: fill one table cell instead
//	}
//
// Response: {"updated": true}.
func handleSetPictureFill(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, _ := v.RequireInt(p, "shape_id")
	if v.HasErrors() {
		return nil, v.Error()
	}

	var params struct {
		ImagePath string `json:"image_path"`
		ImageData string `json:"image_data"`
		Mode      string `json:"mode"`
		Crop      *struct {
			Left   float64 `json:"left"`
			Top    float64 `json:"top"`
			Right  float64 `json:"right"`
			Bottom float64 `json:"bottom"`
		} `json:"crop"`
		Row    *int `json:"row"`
		Column *int `json:"column"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	request := PictureFillRequest{
		ImagePath: params.ImagePath,
		ImageData: params.ImageData,
		Mode:      params.Mode,
		Row:       params.Row,
		Column:    params.Column,
	}
	if params.Crop != nil {
		request.Crop = &PictureFillCropRequest{
			Left:   params.Crop.Left,
			Top:    params.Crop.Top,
			Right:  params.Crop.Right,
			Bottom: params.Crop.Bottom,
		}
	}

	if err := e.SetShapePictureFill(slideIndex, shapeID, request); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return slidescatalog.BuildUpdatedResponse(), nil
}

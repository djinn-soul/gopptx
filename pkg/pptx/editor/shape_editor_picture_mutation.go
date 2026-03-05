package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	minImageCropFraction = 0.0
	maxImageCropFraction = 1.0
	minImageRotationDeg  = -360.0
	maxImageRotationDeg  = 360.0
)

var (
	pictureBlipFillPattern = regexp.MustCompile(`(?s)<p:blipFill\b[^>]*>(.*?)</p:blipFill>`)
	pictureSrcRectPattern  = regexp.MustCompile(`(?s)<a:srcRect\b[^>]*/>`)
	pictureXfrmPattern     = regexp.MustCompile(`<a:xfrm\b([^>]*)>`)
)

func hasPictureUpdateFields(updates common.ShapeUpdate) bool {
	return updates.Crop != nil || updates.Rotation != nil || updates.FlipH != nil || updates.FlipV != nil
}

func validatePictureUpdateFields(updates common.ShapeUpdate) error {
	if updates.Crop != nil {
		crop := updates.Crop
		fields := []struct {
			name  string
			value float64
		}{
			{name: "crop.left", value: crop.Left},
			{name: "crop.right", value: crop.Right},
			{name: "crop.top", value: crop.Top},
			{name: "crop.bottom", value: crop.Bottom},
		}
		for _, field := range fields {
			if field.value < minImageCropFraction || field.value > maxImageCropFraction {
				return fmt.Errorf(
					"%s must be between %.1f and %.1f",
					field.name,
					minImageCropFraction,
					maxImageCropFraction,
				)
			}
		}
	}

	if updates.Rotation != nil {
		rotation := *updates.Rotation
		if rotation < minImageRotationDeg || rotation > maxImageRotationDeg {
			return fmt.Errorf(
				"rotation must be between %.0f and %.0f degrees",
				minImageRotationDeg,
				maxImageRotationDeg,
			)
		}
	}
	return nil
}

func applyPictureShapeUpdates(xmlData []byte, updates common.ShapeUpdate) ([]byte, error) {
	if err := validatePictureUpdateFields(updates); err != nil {
		return nil, err
	}
	updated := xmlData
	var err error

	if updates.Crop != nil {
		updated, err = replacePictureCrop(updated, updates.Crop)
		if err != nil {
			return nil, err
		}
	}
	if updates.Rotation != nil || updates.FlipH != nil || updates.FlipV != nil {
		updated, err = replacePictureTransformAttrs(updated, updates)
		if err != nil {
			return nil, err
		}
	}
	return updated, nil
}

func replacePictureCrop(xmlData []byte, crop *common.ImageCrop) ([]byte, error) {
	match := pictureBlipFillPattern.FindSubmatchIndex(xmlData)
	if match == nil {
		return nil, errors.New("picture shape missing blipFill")
	}

	inner := string(xmlData[match[2]:match[3]])
	inner = pictureSrcRectPattern.ReplaceAllString(inner, "")
	inner = strings.TrimSpace(inner)
	srcRect := buildImageCropXML(&common.ShapeUpdate{Crop: crop})
	inner = strings.TrimSpace(srcRect + inner)

	replacement := "<p:blipFill>" + inner + "</p:blipFill>"
	out := append([]byte{}, xmlData[:match[0]]...)
	out = append(out, replacement...)
	out = append(out, xmlData[match[1]:]...)
	return out, nil
}

func replacePictureTransformAttrs(xmlData []byte, updates common.ShapeUpdate) ([]byte, error) {
	match := pictureXfrmPattern.FindSubmatchIndex(xmlData)
	if match == nil {
		return nil, errors.New("picture shape missing xfrm")
	}

	attrs := string(xmlData[match[2]:match[3]])
	if updates.Rotation != nil {
		attrs = setOrRemoveXMLAttr(attrs, "rot", strconv.Itoa(int(*updates.Rotation*imageRotationScale)), true)
	}
	if updates.FlipH != nil {
		attrs = setOrRemoveXMLAttr(attrs, "flipH", "1", *updates.FlipH)
	}
	if updates.FlipV != nil {
		attrs = setOrRemoveXMLAttr(attrs, "flipV", "1", *updates.FlipV)
	}

	replacement := "<a:xfrm" + attrs + ">"
	out := append([]byte{}, xmlData[:match[0]]...)
	out = append(out, replacement...)
	out = append(out, xmlData[match[1]:]...)
	return out, nil
}

func setOrRemoveXMLAttr(attrs, name, value string, enabled bool) string {
	attrRe := regexp.MustCompile(`\s+` + regexp.QuoteMeta(name) + `="[^"]*"`)
	updated := attrRe.ReplaceAllString(attrs, "")
	if enabled {
		return updated + ` ` + name + `="` + value + `"`
	}
	return updated
}

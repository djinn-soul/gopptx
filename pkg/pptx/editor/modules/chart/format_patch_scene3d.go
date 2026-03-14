package chart

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reChartSpaceSpPr = regexp.MustCompile(`(?s)<c:spPr>.*?</c:spPr>`)
	reScene3DPatch   = regexp.MustCompile(`(?s)<a:scene3d\b.*?</a:scene3d>`)
)

func validateScene3DUpdate(req common.ChartFormatUpdate) error {
	hasSceneField := req.CameraPreset != nil ||
		req.CameraFieldOfView != nil ||
		req.LightRig != nil ||
		req.LightDirection != nil ||
		req.LightRigRevolution != nil
	if !hasSceneField {
		return nil
	}
	if req.CameraPreset == nil || strings.TrimSpace(*req.CameraPreset) == "" {
		return errors.New("camera_preset is required when updating scene3d")
	}
	if req.LightRig == nil || strings.TrimSpace(*req.LightRig) == "" {
		return errors.New("light_rig is required when updating scene3d")
	}
	if req.LightDirection == nil || !isLightDirection(strings.TrimSpace(*req.LightDirection)) {
		return errors.New("light_direction must be one of t,tr,r,br,b,bl,l,tl")
	}
	if req.CameraFieldOfView != nil && *req.CameraFieldOfView <= 0 {
		return errors.New("camera_field_of_view must be > 0")
	}
	return nil
}

func patchChartScene3D(xml string, req common.ChartFormatUpdate) string {
	if req.CameraPreset == nil && req.CameraFieldOfView == nil && req.LightRig == nil &&
		req.LightDirection == nil && req.LightRigRevolution == nil {
		return xml
	}

	scene3D := buildScene3DXML(req)
	if match := reChartSpaceSpPr.FindString(xml); match != "" {
		spPr := match
		if reScene3DPatch.MatchString(spPr) {
			spPr = reScene3DPatch.ReplaceAllString(spPr, scene3D)
		} else {
			spPr = strings.Replace(spPr, "</c:spPr>", scene3D+"</c:spPr>", 1)
		}
		return strings.Replace(xml, match, spPr, 1)
	}

	spPr := "<c:spPr>" + scene3D + "</c:spPr>"
	if idx := strings.Index(xml, "<c:chart>"); idx >= 0 {
		return xml[:idx] + spPr + xml[idx:]
	}
	return xml
}

func buildScene3DXML(req common.ChartFormatUpdate) string {
	var b strings.Builder
	b.WriteString("<a:scene3d>")
	b.WriteString(`<a:camera prst="`)
	b.WriteString(xmlEscape(strings.TrimSpace(*req.CameraPreset)))
	b.WriteString(`"`)
	if req.CameraFieldOfView != nil {
		b.WriteString(` fov="`)
		b.WriteString(strconv.Itoa(*req.CameraFieldOfView))
		b.WriteString(`"`)
	}
	b.WriteString("/>")
	b.WriteString(`<a:lightRig rig="`)
	b.WriteString(xmlEscape(strings.TrimSpace(*req.LightRig)))
	b.WriteString(`" dir="`)
	b.WriteString(xmlEscape(strings.TrimSpace(*req.LightDirection)))
	b.WriteString(`"`)
	if req.LightRigRevolution != nil {
		b.WriteString(` rev="`)
		b.WriteString(boolToOneZero(*req.LightRigRevolution))
		b.WriteString(`"`)
	}
	b.WriteString("/>")
	b.WriteString("</a:scene3d>")
	return b.String()
}

func isLightDirection(direction string) bool {
	switch direction {
	case "t", "tr", "r", "br", "b", "bl", "l", "tl":
		return true
	default:
		return false
	}
}

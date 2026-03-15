package shape

import (
	"errors"
	"fmt"
	"math"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const effectListEmptyXML = `<a:effectLst/>`

func RenderEffectsXML(
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) (string, error) {
	if shadow == nil && glow == nil && blur == nil && softEdge == nil && reflection == nil {
		return "", nil
	}
	inheritHandled, inheritXML, err := renderInheritedShadowEffects(shadow, glow, blur, softEdge, reflection)
	if err != nil {
		return "", err
	}
	if inheritHandled {
		return inheritXML, nil
	}

	items, err := renderEffectItems(shadow, glow, blur, softEdge, reflection)
	if err != nil {
		return "", err
	}
	if items.Len() == 0 {
		return "", nil
	}
	return `<a:effectLst>` + items.String() + `</a:effectLst>`, nil
}

func renderInheritedShadowEffects(
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) (bool, string, error) {
	if shadow == nil || shadow.Inherit == nil {
		return false, "", nil
	}
	if shadow.Color != nil || shadow.BlurEmu != nil || shadow.DistanceEmu != nil || shadow.AngleDeg != nil {
		return false, "", errors.New("shadow.inherit cannot be combined with explicit shadow attributes")
	}
	if glow != nil || blur != nil || softEdge != nil || reflection != nil {
		return false, "", errors.New("shadow.inherit cannot be combined with other explicit effects")
	}
	if *shadow.Inherit {
		return true, "", nil
	}
	return true, effectListEmptyXML, nil
}

func renderEffectItems(
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) (strings.Builder, error) {
	var items strings.Builder
	if err := appendEffectXML(&items, shadow, renderShadowXML); err != nil {
		return items, err
	}
	if err := appendEffectXML(&items, glow, renderGlowXML); err != nil {
		return items, err
	}
	if err := appendEffectXML(&items, blur, renderBlurXML); err != nil {
		return items, err
	}
	if err := appendEffectXML(&items, softEdge, renderSoftEdgeXML); err != nil {
		return items, err
	}
	if err := appendEffectXML(&items, reflection, renderReflectionXML); err != nil {
		return items, err
	}
	return items, nil
}

func appendEffectXML[T any](builder *strings.Builder, value *T, render func(*T) (string, error)) error {
	if value == nil {
		return nil
	}
	effectXML, err := render(value)
	if err != nil {
		return err
	}
	builder.WriteString(effectXML)
	return nil
}

func renderShadowXML(shadow *common.ShapeShadow) (string, error) {
	if shadow == nil {
		return "", nil
	}
	color := "000000"
	if shadow.Color != nil {
		normalized, err := NormalizeHexColor(*shadow.Color)
		if err != nil {
			return "", fmt.Errorf("shadow.color: %w", err)
		}
		color = normalized
	}
	blur := 50800
	if shadow.BlurEmu != nil {
		if *shadow.BlurEmu < 0 {
			return "", errors.New("shadow.blur_emu must be >= 0")
		}
		blur = *shadow.BlurEmu
	}
	dist := 38100
	if shadow.DistanceEmu != nil {
		if *shadow.DistanceEmu < 0 {
			return "", errors.New("shadow.distance_emu must be >= 0")
		}
		dist = *shadow.DistanceEmu
	}
	dir := int64(0)
	if shadow.AngleDeg != nil {
		rotation, err := normalizeRotation(*shadow.AngleDeg)
		if err != nil {
			return "", fmt.Errorf("shadow.angle_deg: %w", err)
		}
		dir = rotation
	}
	return fmt.Sprintf(
		`<a:outerShdw blurRad="%d" dist="%d" dir="%d"><a:srgbClr val="%s"/></a:outerShdw>`,
		blur,
		dist,
		dir,
		color,
	), nil
}

func renderGlowXML(glow *common.ShapeGlow) (string, error) {
	if glow == nil {
		return "", nil
	}
	color := "000000"
	if glow.Color != nil {
		normalized, err := NormalizeHexColor(*glow.Color)
		if err != nil {
			return "", fmt.Errorf("glow.color: %w", err)
		}
		color = normalized
	}
	radius := 38100
	if glow.RadiusEmu != nil {
		if *glow.RadiusEmu < 0 {
			return "", errors.New("glow.radius_emu must be >= 0")
		}
		radius = *glow.RadiusEmu
	}
	return fmt.Sprintf(`<a:glow rad="%d"><a:srgbClr val="%s"/></a:glow>`, radius, color), nil
}

func renderBlurXML(blur *common.ShapeBlur) (string, error) {
	if blur == nil {
		return "", nil
	}
	radius := 50800
	if blur.RadiusEmu != nil {
		if *blur.RadiusEmu < 0 {
			return "", errors.New("blur.radius_emu must be >= 0")
		}
		radius = *blur.RadiusEmu
	}
	return fmt.Sprintf(`<a:blur rad="%d"/>`, radius), nil
}

func renderSoftEdgeXML(softEdge *common.ShapeSoftEdge) (string, error) {
	if softEdge == nil {
		return "", nil
	}
	radius := 50800
	if softEdge.RadiusEmu != nil {
		if *softEdge.RadiusEmu < 0 {
			return "", errors.New("soft_edge.radius_emu must be >= 0")
		}
		radius = *softEdge.RadiusEmu
	}
	return fmt.Sprintf(`<a:softEdge rad="%d"/>`, radius), nil
}

func renderReflectionXML(reflection *common.ShapeReflection) (string, error) {
	if reflection == nil {
		return "", nil
	}
	blur := 0
	if reflection.BlurEmu != nil {
		if *reflection.BlurEmu < 0 {
			return "", errors.New("reflection.blur_emu must be >= 0")
		}
		blur = *reflection.BlurEmu
	}
	dist := 0
	if reflection.DistanceEmu != nil {
		if *reflection.DistanceEmu < 0 {
			return "", errors.New("reflection.distance_emu must be >= 0")
		}
		dist = *reflection.DistanceEmu
	}
	return fmt.Sprintf(`<a:reflection blurRad="%d" dist="%d"/>`, blur, dist), nil
}

func normalizeRotation(raw float64) (int64, error) {
	if math.IsNaN(raw) || math.IsInf(raw, 0) {
		return 0, errors.New("rotation must be finite")
	}
	if raw < -360.0 || raw > 360.0 {
		return 0, errors.New("rotation must be between -360 and 360 degrees")
	}
	return int64(math.Round(raw * rotationDegreeToOOXML)), nil
}

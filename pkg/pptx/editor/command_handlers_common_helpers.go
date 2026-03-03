package editor

func requireSlideIndex(
	e *PresentationEditor,
	payload map[string]any,
	v *PayloadValidator,
) (int, bool) {
	slideIndex, ok := v.RequireInt(payload, "slide_index")
	if !ok {
		return 0, false
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return 0, false
	}
	return slideIndex, true
}

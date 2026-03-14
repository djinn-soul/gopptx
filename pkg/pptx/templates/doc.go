// Package templates provides high-level, typed presentation builders.
//
// Template builders generate common deck structures and can be customized using:
//   - Branding presets (`BrandingSpec`)
//   - Per-slide layout overrides (`LayoutOverrides`)
//
// Example:
//
//	tmpl := templates.SimpleTemplate{
//		Title:   "Quarterly Update",
//		Content: "Highlights",
//		LayoutOverrides: templates.LayoutOverrides{
//			1: elements.SlideLayoutTwoColumn,
//		},
//	}
//	slides, err := tmpl.Build()
//	if err != nil {
//		// handle validation error
//	}
//	_ = slides
package templates

# 19 - Missing Masters, Layouts, and Themes Features

Scope: track missing masters/layouts/themes capabilities in `gopptx` from the provided feature audit screenshots.

## Missing in gopptx

| Feature | Status | TODO |
| --- | --- | --- |
| Theme-aware scheme colors | Missing | [ ] Add theme-aware scheme color token support across shape/text APIs and theme application flows. |
| Color manipulation (lighter/darker) | Missing | [ ] Add lighter/darker color transformation controls for theme-aware color operations. |
| Color opacity | Missing | [ ] Add color-opacity controls in theme/color APIs and serialization paths. |
| Color mix | Missing | [ ] Add color-mix/blend helpers in theme/color APIs for algorithmic palette generation. |
| `.thmx` file support | Missing | [ ] Add import/apply support for `.thmx` theme files. |
| Embedded fonts (theme context) | Missing | [ ] Add/verify embedded-font handling in theme workflows and update APIs/docs accordingly. |
| Material Design color palette | Missing | [ ] Add Material Design preset color palette mapping for theme creation helpers. |
| IBM Carbon color palette | Missing | [ ] Add IBM Carbon preset color palette mapping for theme creation helpers. |
| HandoutLayout type | Missing | [ ] Add HandoutLayout type support in slide layout/master APIs and serialization paths. |

## Capability Depth Gap

| Feature | Status | TODO |
| --- | --- | --- |
| Placeholder types coverage | Limited | [ ] Expand placeholder-type coverage where needed and document supported type matrix clearly. |
| Built-in theme presets coverage | Limited | [ ] Expand built-in theme preset catalog and add examples for each preset family. |
| Slide layout enum coverage (6 -> 12 types) | Limited | [ ] Add 6 more slide layout enum types so `gopptx` reaches 12-type coverage and update docs/tests accordingly. |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and update this file if needed.
- [ ] Add regression tests for theme-aware colors, lighter/darker manipulation, color opacity/mix, `.thmx` handling, embedded-font theme paths, and preset palette helpers once implemented.

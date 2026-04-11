# 18 - Missing Images & Media Features

Scope: track missing image/media capabilities in `gopptx` from the provided feature audit screenshots.

## Missing in gopptx

| Feature | Status | TODO |
| --- | --- | --- |
| Image from URL | Missing | [ ] Add direct image-from-URL ingestion helper with fetch and validation controls. |
| Image from bytes | Missing | [ ] Add/verify raw-byte image ingestion APIs and bridge mapping. |
| SVG | Missing | [ ] Add SVG image ingestion/render path in image APIs and bridge payloads. |
| Image crop | Missing | [ ] Add/verify crop controls for image frames with stable coordinate semantics. |
| Image rotation | Missing | [ ] Add/verify image rotation controls in image APIs and bridge payloads. |
| Image opacity / transparency | Missing | [ ] Add image opacity/transparency controls in image styling APIs. |
| Animated GIF | Missing | [ ] Add animated GIF support strategy (preserve animation vs first-frame fallback) and API controls. |
| Image rounding (circular) | Missing | [ ] Add circular/rounded image mask helpers and bridge mapping. |
| Image shadow effect | Missing | [ ] Add image shadow effect controls in image styling APIs. |
| Image reflection effect | Missing | [ ] Add image reflection effect controls and serialization. |
| Image glow effect | Missing | [ ] Add image glow effect controls and serialization. |
| Image soft edges | Missing | [ ] Add soft-edge controls for image effects. |
| Image blur effect | Missing | [ ] Add blur controls for image effects. |
| Image inner shadow | Missing | [ ] Add inner-shadow controls for image effects. |
| Contain / Cover sizing modes | Missing | [ ] Add contain/cover fit-mode options for image placement APIs. |
| Alt text on image | Missing | [ ] Add image alt-text controls and persistence/inspection APIs. |
| YouTube embed | Missing | [ ] Add YouTube/video-link embedding helpers with supported playback behavior constraints. |
| 3D model embedding (GLB/GLTF) | Missing | [ ] Add 3D model embedding support and related media-part handling in APIs/bridge. |

## Verification Follow-up

- [ ] Verify listed gaps against current runtime behavior and update this file if needed.
- [ ] Add regression tests for URL/bytes ingestion, crop/rotation/opacity, image effects, alt text, SVG, animated GIF, rounding/mask, contain-cover fit modes, YouTube embedding, and 3D model embedding once implemented.

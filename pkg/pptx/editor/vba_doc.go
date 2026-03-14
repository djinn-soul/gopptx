package editor

// VBA Macro Integration
//
// PresentationEditor preserves macro-enabled package wiring for .pptm files and
// treats `ppt/vbaProject.bin` as an opaque blob.
//
// Security notes:
//   - Embedded VBA is not inspected, linted, or sanitized by gopptx.
//   - Macro signatures are not created or repaired by gopptx.
//   - Editing package parts can invalidate existing macro trust/signature state.
//   - `pkg/pptx/vba` validates `vbaProject.bin` CFB structure with mscfb but
//     does not modify/recompile VBA streams.
//
// Recommended workflow:
//  1. Open and edit macro-enabled decks only from trusted sources.
//  2. Save macro-enabled outputs with the `.pptm` extension.
//  3. Re-sign VBA projects using enterprise signing tooling after edits.
//  4. Validate generated outputs with OpenXML/package checks and real Office open tests.

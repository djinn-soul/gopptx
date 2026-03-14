package editor

// Theme Management API
//
// Lifecycle warning:
// Theme mutation affects every slide/layout/notes master bound to the edited
// theme part. Use ThemeInventory() before edits to inspect bindings and avoid
// unintended global visual changes.
//
// Recommended flow:
//  1. Open a deck with OpenPresentationEditor.
//  2. Inspect ThemeInventory for theme part ownership.
//  3. Apply one of:
//     - SetThemeData(path, xml)
//     - SetThemeFontScheme(major, minor)
//     - SetThemeColorScheme(...)
//     - SetGlobalThemePreset(name)
//  4. Save to a new output file and validate visual impact.

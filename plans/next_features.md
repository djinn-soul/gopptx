# Task: Planning Next Features

## Current Status
Task 26 (VBA Macros) has been finalized. Validation and review fixes for CustomXML and VBA paths are complete.

## Proposed Next Features

### 1. [COMPLEX] Validation & Repair (Task 20)
- **Goal**: Integrate the `PptxValidator` into the core library or a `gopptx validate` command.
- **Value**: Ensures all generated decks are OOXML compliant and openable in all PowerPoint versions.

### 2. [COMPLEX] Export Path: HTML/PDF (Task 21)
- **Goal**: Add baseline support for exporting slides to HTML or PDF.
- **Value**: High utility for web previews and document sharing.

### 3. [MEDIUM] Media Embedding (Task 23)
- **Goal**: Support embedding Video (`.mp4`) and Audio (`.wav`, `.mp3`) files.
- **Value**: Important for rich media presentations.

### 4. [MEDIUM] SmartArt Breadth (Task 24)
- **Goal**: Expand beyond the initial 25 layouts to full dynamic graph generation for all SmartArt types.
- **Value**: Significantly improved visual capability.

### 5. [SIMPLE] Template Expansion (Task 16)
- **Goal**: Complete the "Data Binding" and "Layout Customization" hooks for existing templates.
- **Value**: Makes the template system much more powerful for automated reporting.

### 6. [MEDIUM] Legacy PPT Interop (Task 52)
- **Goal**: Baseline support for reading legacy `.ppt` files (potentially via bridge or external tool).
- **Value**: Enterprise compatibility.

### 7. [EASY] CLI Polish (Task 15)
- **Goal**: Finalize the shell completion and info-display subcommands.
- **Value**: Better developer experience.

# Specification: Initialize core PPTX structure

## Goal
Establish the foundational data structures and XML serialization logic required to generate a valid, minimal PowerPoint (.pptx) file using Go.

## Requirements
- Define core structs for Presentation, Slide, and SlideLayout.
- Implement XML marshaling for the primary PPTX components (`presentation.xml`, `slide1.xml`, etc.).
- Create a basic "Package" manager to handle the ZIP/OPC structure of the .pptx file.
- Support adding a blank slide to a new presentation.
- Ensure cross-platform compatibility for the generated files (readable by PowerPoint, Keynote, Google Slides).

## Technical Details
- **Language:** Go
- **Packages:** `encoding/xml`, `archive/zip`
- **Output:** A `.pptx` file containing a single slide.

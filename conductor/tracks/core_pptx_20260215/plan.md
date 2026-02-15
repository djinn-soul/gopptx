# Implementation Plan - Initialize core PPTX structure

## Phase 1: Foundational Structure [checkpoint: f09ddff]
- [x] Task: Define core PPTX data models (Presentation, Slide) in Go [370c97c]
    - [x] Write Tests: Define expected struct behavior and XML tags
    - [x] Implement Feature: Create structs in `pkg/gopptx/models.go`
- [x] Task: Implement OPC (Open Packaging Conventions) ZIP wrapper [8d47656]
    - [x] Write Tests: Verify ZIP creation and directory structure
    - [x] Implement Feature: Create package manager in `internal/opc/`
- [x] Task: Conductor - User Manual Verification 'Phase 1: Foundational Structure' (Protocol in workflow.md) [f09ddff]

## Phase 2: XML Serialization & Slide Generation
- [x] Task: Implement XML marshaling for Presentation and Slide [49c910f]
    - [x] Write Tests: Verify generated XML against PPTX schema
    - [x] Implement Feature: Add `MarshalXML` methods to models
- [ ] Task: Implement basic Slide Addition logic
    - [ ] Write Tests: Test `presentation.AddSlide()` functionality
    - [ ] Implement Feature: Logic to create new slide files and update relationships
- [ ] Task: Create a smoke test to generate a valid .pptx file
    - [ ] Write Tests: Generate file and verify it can be opened (basic check)
    - [ ] Implement Feature: `examples/basic_gen/main.go`
- [ ] Task: Conductor - User Manual Verification 'Phase 2: XML Serialization & Slide Generation' (Protocol in workflow.md)

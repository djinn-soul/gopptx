# Plan: Task 18 - Layout Helpers

Implement layout helper functions in `gopptx` to simplify positioning and distribution of elements.

## Purpose
Provide developers with idiomatic Go tools to handle common layout tasks like centering, grid alignment, and stacking, reducing manual EMU math.

## Progress
- [ ] Define slide dimension constants in `pkg/pptx/layout_helpers.go` <!-- id: 0 -->
- [ ] Implement `Center(cx, cy int64) (x, y int64)` helper <!-- id: 1 -->
- [ ] Implement `Grid(rows, cols int, margin int64) ([]Box, error)` helper <!-- id: 2 -->
- [ ] Implement `Stack(orientation string, gap int64, elements ...Size) ([]Point, error)` helper <!-- id: 3 -->
- [ ] Implement `Distribute(orientation string, bounds Box, count int) ([]int64, error)` helper <!-- id: 4 -->
- [ ] Add unit tests in `pkg/pptx/layout_helpers_test.go` <!-- id: 5 -->
- [ ] Add GoDoc documentation with EMU conversion examples <!-- id: 6 -->

## Decisions
- **Dimensions**: Standard 4:3 slide (9144000 x 6858000 EMU) as default.
- **Return Types**: Use simple structs or slices of coordinates to keep it lightweight.
- **Error Handling**: Fail-fast if input parameters (like rows/cols <= 0) are invalid.

## Validation
- All unit tests pass with deterministic results.
- GoDoc examples correctly demonstrate EMU usage.

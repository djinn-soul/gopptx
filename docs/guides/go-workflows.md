# Go Workflows

Use this page for practical Go usage patterns.

## 1) Create a Deck

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	p := pptx.NewPresentation()
	s := p.AddSlide()
	s.AddTextBox("Hello from Go")
	_ = p.Save("go_created.pptx")
}
```

## 2) Run Example Modules

```bash
go run ./examples/01-basic-pptx-generation/basic_gen.go
go run ./examples/12-shapes/feature_showcase.go
go run ./examples/58-gopptx-rich-slide/main.go
```

## 3) Edit Pipelines

Use `examples/19-read-modify-existing/` and `examples/43-presentation-props-editor/` as references for read-modify-save workflows.

## Related

- [Go Library](go-library.md)
- [Examples Map](../examples-map.md)

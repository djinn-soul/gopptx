# Go Library

## Install

```bash
go get github.com/djinn-soul/gopptx
```

## Minimal Example

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	p := pptx.NewPresentation()
	s := p.AddSlide()
	s.AddTextBox("Hello from Go")
	_ = p.Save("output.pptx")
}
```

## Practical Usage Notes

- Keep operations in-memory and save once per workflow step.
- For cross-language integrations, use bridge operation contracts rather than exposing internal structs.
- Use `examples/` to validate capability coverage quickly.

## Example Entry Points

Examples are organized by task in `examples/`:

- `01-basic-pptx-generation`
- `19-read-modify-existing`
- `23-media-embed`
- `37-slide-duplication`
- `45-commenting-api`
- `57-placeholder-overrides`

See [Examples Map](../examples-map.md) for runnable commands.

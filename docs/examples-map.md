# Examples Map

The repository has task-focused examples under `examples/`.

## Recommended Starting Points

- `01-basic-pptx-generation`: first deck generation
- `04-text-formatting`: title/body text controls
- `09-charts`: chart families and chart data updates
- `12-shapes`: shape APIs and layout helpers
- `19-read-modify-existing`: edit existing decks
- `43-presentation-props-editor`: metadata and theme/reskin workflows
- `58-gopptx-rich-slide`: rich all-in-one feature sample

## Example Output Layout

- Source code: `examples/<example-name>/`
- Generated PPTX files: `examples/output/*.pptx`
- Rendered screenshots: `examples/output/png/**/Slide*.PNG`

## Run Any Example

```bash
go run ./examples/01-basic-pptx-generation/basic_gen.go
```

```bash
go run ./examples/09-charts/chart_smoke.go
```

## Docs Gallery Cross-Links

- [Basic Generation](showcase/basic-generation.md)
- [Basic Usage](showcase/basic-usage.md)
- [Rich Slide](showcase/rich-slide.md)
- [Text Styling](showcase/text-styling.md)
- [Radar Chart](showcase/chart-radar.md)

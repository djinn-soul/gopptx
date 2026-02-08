# pptcli

`pptcli` is the command-line entrypoint for generating and validating PPTX files in `gopptx`.

## Commands

### `create`

Create a new presentation with generated slide titles.

```bash
go run ./cmd/pptcli create -out demo.pptx -title "Demo Deck" -slides 3
```

### `md2ppt`

Convert markdown into PPTX slides.

```bash
go run ./cmd/pptcli md2ppt -in deck.md -out deck.pptx -title "Deck From Markdown"
```

If `-out` is omitted, output defaults to `<input-stem>.pptx` next to the input markdown file.

### `info`

Print metadata and package summary information for a PPTX file.

```bash
go run ./cmd/pptcli info -file deck.pptx
```

### `validate`

Validate that a file is a readable PPTX package with required parts and well-formed XML.

```bash
go run ./cmd/pptcli validate -file deck.pptx
```

## Legacy flags mode

The historical flag-only entrypoint is still supported:

```bash
go run ./cmd/pptcli -out output.pptx
go run ./cmd/pptcli -md deck.md -title "Deck Title" -out output.pptx
```




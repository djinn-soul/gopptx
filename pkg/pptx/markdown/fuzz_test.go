package markdown

import "testing"

func FuzzSlidesFromMarkdown(f *testing.F) {
	f.Add("# Title\n\n- bullet\n- another\n\n---\n\n# Slide 2\n\nParagraph text.")
	f.Add("# Hello\n\n**bold** and *italic* text")
	f.Add("# Code slide\n\n```go\nfmt.Println(\"hi\")\n```")
	f.Add("| col1 | col2 |\n|---|---|\n| a | b |")
	f.Add("---")
	f.Add("")
	f.Add("1. first\n2. second\n3. third")
	f.Add("# Img\n\n![alt](image.png)")
	f.Fuzz(func(_ *testing.T, md string) {
		_, _ = SlidesFromMarkdown(md)
	})
}

func FuzzParseInlineTextRuns(f *testing.F) {
	f.Add("plain text")
	f.Add("**bold** and *italic*")
	f.Add("`code` inline")
	f.Add("**bold *nested italic* still bold**")
	f.Add("")
	f.Add("***all***")
	f.Add("**")
	f.Add("*")
	f.Add("`")
	f.Fuzz(func(_ *testing.T, text string) {
		_, _ = parseInlineTextRuns(text)
	})
}

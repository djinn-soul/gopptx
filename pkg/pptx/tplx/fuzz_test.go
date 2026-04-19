package tplx

import "testing"

func FuzzMergeAdjacentRuns(f *testing.F) {
	f.Add([]byte(`<?xml version="1.0" encoding="UTF-8"?><p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><p:cSld><p:spTree><p:sp><p:txBody><a:p><a:r><a:t>hello </a:t></a:r><a:r><a:t>world</a:t></a:r></a:p></p:txBody></p:sp></p:spTree></p:cSld></p:sld>`))
	f.Add([]byte(`<a:p xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><a:r><a:t>{{na</a:t></a:r><a:r><a:t>me}}</a:t></a:r></a:p>`))
	f.Add([]byte("<a:p></a:p>"))
	f.Add([]byte(""))
	f.Add([]byte("not xml at all"))
	f.Add([]byte(`<a:p><a:r><a:rPr/><a:t>A</a:t></a:r><a:r><a:rPr/><a:t>B</a:t></a:r></a:p>`))
	f.Fuzz(func(_ *testing.T, xmlBytes []byte) {
		_ = mergeAdjacentRuns(xmlBytes)
	})
}

func FuzzInterpolateText(f *testing.F) {
	f.Add("hello {{name}}", "name", "World")
	f.Add("{{greeting}} {{name}}!", "greeting", "Hi")
	f.Add("no placeholders here", "key", "value")
	f.Add("{{missing}}", "key", "value")
	f.Add("", "key", "value")
	f.Add("{{a}}{{b}}{{c}}", "a", "1")
	f.Add("{{ key with spaces }}", "key", "val")
	f.Add("{{<script>}}", "key", "val")
	f.Fuzz(func(_ *testing.T, text, ctxKey, ctxVal string) {
		ctx := Context{ctxKey: ctxVal}
		_, _ = interpolateText(text, ctx, nil, false)
	})
}

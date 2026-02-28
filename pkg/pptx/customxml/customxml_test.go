package customxml

import (
	"reflect"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func TestCustomXmlStoreBasic(t *testing.T) {
	s := NewStore()
	if !s.IsEmpty() {
		t.Fatalf("expected store to be empty initially")
	}

	p := s.Add("MyRoot")
	p.Namespace("http://example.com/ns")
	p.Property("foo", "bar")
	p.Property("baz", "qux")
	p.Content("<child>data</child>")

	if s.IsEmpty() || s.Len() != 1 {
		t.Fatalf("expected store to have 1 item")
	}

	parts := s.ToCommonParts()
	if len(parts) != 1 {
		t.Fatalf("expected 1 converted part")
	}

	got := parts[0]
	want := common.CustomXMLPart{
		RootElement: "MyRoot",
		Namespace:   "http://example.com/ns",
		Content:     "<child>data</child>",
		Properties: []common.CustomXMLKV{
			{Key: "foo", Value: "bar"},
			{Key: "baz", Value: "qux"},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got part %+v; want %+v", got, want)
	}
}

func TestCustomXmlStoreRaw(t *testing.T) {
	s := NewStore()
	s.AddRaw(common.CustomXMLPart{Content: "<raw/>"})

	if s.Len() != 1 {
		t.Fatalf("expected store to have 1 item")
	}
	parts := s.ToCommonParts()
	if parts[0].Content != "<raw/>" {
		t.Errorf("expected <raw/> content")
	}
}

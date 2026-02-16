package common_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestResolveRelationshipTarget(t *testing.T) {
	tests := []struct {
		source string
		target string
		want   string
	}{
		{"ppt/slides/slide1.xml", "../media/image1.png", "ppt/media/image1.png"},
		{"ppt/slides/slide1.xml", "../layouts/layout1.xml", "ppt/layouts/layout1.xml"},
		{"ppt/presentation.xml", "slides/slide1.xml", "ppt/slides/slide1.xml"},
		{"ppt/charts/chart1.xml", "../embeddings/sheet1.xlsx", "ppt/embeddings/sheet1.xlsx"},
	}

	for _, tt := range tests {
		got := common.ResolveRelationshipTarget(tt.source, tt.target)
		if got != tt.want {
			t.Errorf("ResolveRelationshipTarget(%q, %q) = %q; want %q", tt.source, tt.target, got, tt.want)
		}
	}
}

func TestMakeRelativePath(t *testing.T) {
	tests := []struct {
		source string
		target string
		want   string
	}{
		{"ppt/slides/slide1.xml", "ppt/media/image1.png", "../media/image1.png"},
		{"ppt/slides/slide1.xml", "ppt/layouts/layout1.xml", "../layouts/layout1.xml"},
		{"ppt/presentation.xml", "ppt/slides/slide1.xml", "slides/slide1.xml"},
		{"ppt/charts/chart1.xml", "ppt/embeddings/sheet1.xlsx", "../embeddings/sheet1.xlsx"},
	}

	for _, tt := range tests {
		got := common.MakeRelativePath(tt.source, tt.target)
		if got != tt.want {
			t.Errorf("MakeRelativePath(%q, %q) = %q; want %q", tt.source, tt.target, got, tt.want)
		}
	}
}

package pptx_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestSlideBuilderIsAvailableFromPrimaryFacade(t *testing.T) {
	builder := pptx.NewSlideBuilder("Facade")
	requireSlideBuilder(builder)
	slide := builder.AddBullet("one").Build()
	rebuilt := pptx.BuildFrom(slide).AddBullet("two").Build()

	if rebuilt.Title != "Facade" || len(rebuilt.Bullets) != 2 {
		t.Fatalf("unexpected facade-built slide: %#v", rebuilt)
	}
}

func requireSlideBuilder(_ *pptx.SlideBuilder) {}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vegito/goppt/pkg/pptx"
)

func main() {
	var out string
	flag.StringVar(&out, "out", "output.pptx", "output PPTX file")
	flag.Parse()

	slides := []pptx.SlideContent{
		pptx.NewSlide("Welcome").AddBullet("Ported from ppt-rs").AddBullet("Go baseline ready"),
		pptx.NewSlide("Next Steps").AddBullet("Add markdown parser").AddBullet("Add tables/charts/images"),
	}

	if err := pptx.WriteFile(out, "Go PPTX", slides); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

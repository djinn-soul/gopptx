package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vegito/goppt/pkg/pptx"
)

func main() {
	var out string
	var markdownPath string
	var title string
	flag.StringVar(&out, "out", "output.pptx", "output PPTX file")
	flag.StringVar(&markdownPath, "md", "", "markdown file to convert to PPTX")
	flag.StringVar(&title, "title", "", "presentation title (required with -md)")
	flag.Parse()

	if markdownPath != "" {
		if title == "" {
			fmt.Fprintln(os.Stderr, "title is required when using -md")
			os.Exit(1)
		}

		markdown, err := os.ReadFile(markdownPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		slides, err := pptx.SlidesFromMarkdown(string(markdown))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := pptx.WriteFile(out, title, slides); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Welcome").AddBullet("Ported from ppt-rs").AddBullet("Go baseline ready"),
		pptx.NewSlide("Next Steps").AddBullet("Add markdown parser").AddBullet("Add tables/charts/images"),
	}

	if err := pptx.WriteFile(out, "Go PPTX", slides); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

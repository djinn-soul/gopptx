package main

import (
	"log"

	"github.com/djinn-soul/gopptx/pkg/gopptx"
)

func main() {
	pres := &gopptx.Presentation{}
	pres.AddSlide()

	filename := "hello_world.pptx"
	err := pres.Save(filename)
	if err != nil {
		log.Fatalf("Failed to save presentation: %v", err)
	}

	log.Printf("Successfully generated %s\n", filename)
}

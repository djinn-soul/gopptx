package gopptx

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	NSP = "http://schemas.openxmlformats.org/presentationml/2006/main"
)

// Presentation represents the main presentation XML component.
type Presentation struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main presentation"`
	Slides  []*Slide `xml:"sldIdLst>sldId"`
}

// AddSlide adds a new blank slide to the presentation.
func (p *Presentation) AddSlide() *Slide {
	slide := &Slide{}
	p.Slides = append(p.Slides, slide)
	return slide
}

// Save writes the presentation to a .pptx file.
func (p *Presentation) Save(path string) error {
	slideCount := len(p.Slides)
	if slideCount < 1 {
		return errors.New("at least one slide is required")
	}

	data, err := pptx.Create("Presentation", slideCount)
	if err != nil {
		return err
	}
	return writeFileAtomically(path, data, 0o600)
}

// Slide represents a single slide XML component.
type Slide struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
}

func writeFileAtomically(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".gopptx-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	committed := false
	defer func() {
		_ = tmpFile.Close()
		if !committed {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err = tmpFile.Write(content); err != nil {
		return err
	}
	if err = tmpFile.Chmod(perm); err != nil {
		return err
	}
	if err = tmpFile.Sync(); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	if _, err = os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return err
	}
	committed = true
	return nil
}

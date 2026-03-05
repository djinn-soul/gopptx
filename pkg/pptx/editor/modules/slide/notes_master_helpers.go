package slide

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type ReadFileFn func(string) ([]byte, error)
type RegisterImageFn func([]byte, string) (string, error)

func ResolveNotesMasterBackgroundMedia(
	master *elements.NotesMaster,
	readFile ReadFileFn,
	registerImage RegisterImageFn,
) (string, []string, error) {
	if master == nil || master.Background == nil ||
		master.Background.Type != elements.SlideBackgroundPicture || master.Background.PictureFill == nil {
		return "", nil, nil
	}

	img := master.Background.PictureFill
	var data []byte
	var err error

	if len(img.Data) > 0 {
		data = img.Data
	} else if img.Path != "" {
		data, err = readFile(img.Path)
		if err != nil {
			return "", nil, fmt.Errorf("read background image: %w", err)
		}
	}

	if len(data) == 0 {
		return "", nil, nil
	}

	format := img.Format
	if format == "" {
		if img.Path != "" {
			format = strings.TrimPrefix(filepath.Ext(img.Path), ".")
		}
		if format == "" {
			format = "png"
		}
	}

	internalPath, err := registerImage(data, format)
	if err != nil {
		return "", nil, fmt.Errorf("register background image: %w", err)
	}

	target := "../media/" + path.Base(internalPath)
	return "rId2", []string{target}, nil
}

package presentation

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

//nolint:gocognit // Per-slide comment author/index normalization is intentionally explicit.
func prepareComments(
	meta Metadata,
	slides []elements.SlideContent,
) ([]comments.Author, map[int][]comments.Comment, []int) {
	var authors []comments.Author
	authorMap := make(map[string]int)

	slideComments := make(map[int][]comments.Comment)
	var commentSlideIndices []int

	commentDate := meta.GeneratedDate
	if commentDate.IsZero() {
		commentDate = time.Now()
	}

	for i, s := range slides {
		if len(s.Comments) == 0 {
			continue
		}
		var slideCms []comments.Comment
		for _, c := range s.Comments {
			idx, ok := authorMap[c.AuthorName]
			if !ok {
				idx = len(authors)
				authorMap[c.AuthorName] = idx

				initials := ""
				parts := strings.FieldsSeq(c.AuthorName)
				var initialsSb290 strings.Builder
				for p := range parts {
					if r, _ := utf8.DecodeRuneInString(p); r != utf8.RuneError {
						initialsSb290.WriteRune(r)
					}
				}
				initials += initialsSb290.String()
				if len([]rune(initials)) > maxAuthorInitialRunes {
					initials = string([]rune(initials)[:maxAuthorInitialRunes])
				}
				if initials == "" {
					initials = "A"
				}

				authors = append(authors, comments.Author{
					ID:         int64(idx + 1),
					Name:       c.AuthorName,
					Initials:   initials,
					ColorIndex: idx % authorColorPaletteSize,
					LastIndex:  0,
				})
			}

			authors[idx].LastIndex++
			x := c.X
			y := c.Y
			if x == 0 && y == 0 {
				x = 100000
				y = 100000
			}

			slideCms = append(slideCms, comments.Comment{
				ID:       int64(len(slideCms) + 1),
				AuthorID: authors[idx].ID,
				Text:     c.Text,
				Date:     commentDate,
				X:        x,
				Y:        y,
				Index:    authors[idx].LastIndex,
			})
		}
		slideComments[i] = slideCms
		commentSlideIndices = append(commentSlideIndices, i+1)
	}
	return authors, slideComments, commentSlideIndices
}

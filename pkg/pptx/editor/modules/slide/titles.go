package slide

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"sort"
	"strings"
	"sync"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func PopulateSlideTitlesConcurrently(
	slides []common.EditorSlideRef,
	readPart func(string) ([]byte, bool),
) []string {
	if len(slides) == 0 {
		return nil
	}

	type result struct {
		index int
		title string
	}
	ch := make(chan result, len(slides))
	var wg sync.WaitGroup

	for idx := range slides {
		wg.Go(func() {
			data, _ := readPart(slides[idx].Part)
			title := ExtractFirstAText(data)
			ch <- result{index: idx, title: title}
		})
	}
	wg.Wait()
	close(ch)

	results := make([]result, 0, len(slides))
	for item := range ch {
		results = append(results, item)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].index < results[j].index })

	titles := make([]string, len(slides))
	for _, item := range results {
		titles[item.index] = item.title
	}
	return titles
}

func ExtractFirstAText(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			return ""
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "t" {
			continue
		}
		var value string
		if decodeErr := decoder.DecodeElement(&value, &start); decodeErr != nil {
			return ""
		}
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
}

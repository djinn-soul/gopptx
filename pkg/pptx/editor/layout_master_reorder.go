package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	slideLayoutListOpen  = "<p:sldLayoutIdLst>"
	slideLayoutListClose = "</p:sldLayoutIdLst>"
	firstSlideLayoutID   = 2147483649
)

var slideLayoutRelIDPattern = regexp.MustCompile(`<p:sldLayoutId[^>]*r:id="([^"]+)"`)

// ReorderSlideLayouts reorders layout references within a slide master's <p:sldLayoutIdLst>.
func (e *PresentationEditor) ReorderSlideLayouts(masterPart string, layoutParts []string) error {
	masterPart = common.CanonicalPartPath(masterPart)
	if masterPart == "" {
		masters, err := e.ListSlideMasters()
		if err != nil || len(masters) == 0 {
			return errors.New("no slide masters found in presentation")
		}
		masterPart = common.CanonicalPartPath(masters[0].Part)
	}

	xmlData, ok := e.parts.Get(masterPart)
	if !ok {
		return fmt.Errorf("master part %s not found", masterPart)
	}
	pathToRelID, err := e.layoutRelationshipIDs(masterPart)
	if err != nil {
		return err
	}

	content := string(xmlData)
	beforeList, listAndAfter, foundOpen := strings.Cut(content, slideLayoutListOpen)
	currentList, afterList, foundClose := strings.Cut(listAndAfter, slideLayoutListClose)
	if !foundOpen || !foundClose {
		return nil
	}

	orderedRelIDs := reorderedLayoutRelIDs(currentList, layoutParts, pathToRelID)
	var newList strings.Builder
	newList.WriteString(slideLayoutListOpen)
	for index, relID := range orderedRelIDs {
		_, _ = fmt.Fprintf(
			&newList,
			`<p:sldLayoutId id="%d" r:id="%s"/>`,
			firstSlideLayoutID+index,
			relID,
		)
	}
	newList.WriteString(slideLayoutListClose)
	e.parts.Set(masterPart, []byte(beforeList+newList.String()+afterList))
	return nil
}

func (e *PresentationEditor) layoutRelationshipIDs(masterPart string) (map[string]string, error) {
	relsData, _ := e.parts.Get(common.RelsPathFor(masterPart))
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil, fmt.Errorf("parse master rels: %w", err)
	}

	pathToRelID := make(map[string]string)
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideLayout {
			targetPart := common.CanonicalPartPath(common.ResolveRelationshipTarget(masterPart, rel.Target))
			pathToRelID[targetPart] = rel.ID
		}
	}
	return pathToRelID, nil
}

func reorderedLayoutRelIDs(
	currentList string,
	layoutParts []string,
	pathToRelID map[string]string,
) []string {
	ordered := make([]string, 0, len(pathToRelID))
	seen := make(map[string]bool, len(pathToRelID))
	for _, layoutPart := range layoutParts {
		relID, found := pathToRelID[common.CanonicalPartPath(layoutPart)]
		if found && !seen[relID] {
			ordered = append(ordered, relID)
			seen[relID] = true
		}
	}
	for _, match := range slideLayoutRelIDPattern.FindAllStringSubmatch(currentList, -1) {
		relID := match[1]
		if !seen[relID] {
			ordered = append(ordered, relID)
			seen[relID] = true
		}
	}
	return ordered
}

package layoutmaster

import (
	"fmt"
	"regexp"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

var masterNumExtractPattern = regexp.MustCompile(`slideMaster(\d+)\.xml`)

func AddMasterRelationship(
	nonSlideRels []common.EditorRelationship,
	presentationXML string,
	nextRelIDNum int,
	masterPart string,
) ([]common.EditorRelationship, string, int, error) {
	newMasterRelID := fmt.Sprintf("rId%d", nextRelIDNum)
	nextRelIDNum++

	nonSlideRels = append(nonSlideRels, common.EditorRelationship{
		ID:     newMasterRelID,
		Type:   common.RelTypeSlideMaster,
		Target: common.MakeRelativePath(common.PresentationXMLPath, masterPart),
	})

	updatedPresentationXML, err := editorslide.RewritePresentationSlideMasterList(
		[]byte(presentationXML),
		newMasterRelID,
	)
	if err != nil {
		return nil, "", 0, err
	}

	return nonSlideRels, updatedPresentationXML, nextRelIDNum, nil
}

func NextLayoutRelID(masterRels []common.EditorRelationship) string {
	maxID := 0
	for _, rel := range masterRels {
		id, ok := common.ParseRelationshipNumber(rel.ID)
		if ok && id > maxID {
			maxID = id
		}
	}
	return fmt.Sprintf("rId%d", maxID+1)
}

func AppendLayoutRelationship(
	masterRels []common.EditorRelationship,
	masterPart string,
	layoutPart string,
) []common.EditorRelationship {
	newRelID := NextLayoutRelID(masterRels)
	return append(masterRels, common.EditorRelationship{
		ID:     newRelID,
		Type:   common.RelTypeSlideLayout,
		Target: common.MakeRelativePath(masterPart, layoutPart),
	})
}

func FilterOutRelationshipTarget(
	rels []common.EditorRelationship,
	target string,
) []common.EditorRelationship {
	filtered := make([]common.EditorRelationship, 0, len(rels))
	for _, rel := range rels {
		if rel.Target == target {
			continue
		}
		filtered = append(filtered, rel)
	}
	return filtered
}

func ExtractMasterNumber(masterPart string) int {
	matches := masterNumExtractPattern.FindStringSubmatch(masterPart)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 1
}

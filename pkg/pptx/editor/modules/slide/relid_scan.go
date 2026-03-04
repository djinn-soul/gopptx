package slide

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

func NextRelationshipIDNum(slides []common.EditorSlideRef, nonSlideRels []common.EditorRelationship) int {
	maxNum := 0
	for _, slide := range slides {
		if num, ok := common.ParseRelationshipNumber(slide.RelID); ok && num > maxNum {
			maxNum = num
		}
	}
	for _, rel := range nonSlideRels {
		if num, ok := common.ParseRelationshipNumber(rel.ID); ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

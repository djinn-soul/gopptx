package slide

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

func NextSlideID(slides []common.EditorSlideRef) int64 {
	var maxID int64 = 255
	for _, slide := range slides {
		if slide.SlideID > maxID {
			maxID = slide.SlideID
		}
	}
	return maxID + 1
}

func NextSlidePartNumber(slides []common.EditorSlideRef) int {
	maxNum := 0
	for _, slide := range slides {
		num, ok := common.ParseSlidePartNumber(slide.Part)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

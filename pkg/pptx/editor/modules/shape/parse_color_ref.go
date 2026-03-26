package shape

func parseColorRef(src *struct {
	SrgbClr *struct {
		Val string `xml:"val,attr"`
	} `xml:"srgbClr"`
}) (string, bool) {
	if src == nil || src.SrgbClr == nil || src.SrgbClr.Val == "" {
		return "", false
	}
	return src.SrgbClr.Val, true
}

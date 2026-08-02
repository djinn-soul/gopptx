package editorcommon

// GroupChildSpace is a group shape's child coordinate space, the a:chOff and
// a:chExt of its <a:xfrm>.
//
// A child inside a group states its position in this space, not in slide EMUs.
// PowerPoint maps one to the other by scaling the child extent onto the group
// extent, so a group whose chExt differs from its ext draws its children
// somewhere other than where their raw numbers say.
type GroupChildSpace struct {
	OffsetX  int `json:"offset_x"`
	OffsetY  int `json:"offset_y"`
	ExtentCx int `json:"extent_cx"`
	ExtentCy int `json:"extent_cy"`
}

// ChildToSlide maps a child's coordinates from the group's child space onto the
// slide, given where the group itself sits (groupX, groupY) and how big it is
// (groupW, groupH).
//
// A child space with a zero extent carries no scale -- PowerPoint treats the
// two spaces as identical then -- so the child is only translated.
func (s *GroupChildSpace) ChildToSlide(
	groupX, groupY, groupW, groupH int,
	childX, childY, childW, childH int,
) (int, int, int, int) {
	if s == nil {
		return childX, childY, childW, childH
	}
	scaleX, scaleY := 1.0, 1.0
	if s.ExtentCx > 0 && groupW > 0 {
		scaleX = float64(groupW) / float64(s.ExtentCx)
	}
	if s.ExtentCy > 0 && groupH > 0 {
		scaleY = float64(groupH) / float64(s.ExtentCy)
	}
	x := groupX + int(float64(childX-s.OffsetX)*scaleX)
	y := groupY + int(float64(childY-s.OffsetY)*scaleY)
	w := int(float64(childW) * scaleX)
	h := int(float64(childH) * scaleY)
	return x, y, w, h
}

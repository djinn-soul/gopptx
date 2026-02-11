package common

// Point represents a 2D coordinate in EMU.
type Point struct {
	X int64
	Y int64
}

// Size represents dimensions in EMU.
type Size struct {
	CX int64
	CY int64
}

// Box represents a rectangular region with position and size in EMU.
type Box struct {
	X  int64
	Y  int64
	CX int64
	CY int64
}

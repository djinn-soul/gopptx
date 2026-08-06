package mermaid

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

const (
	// relationLabelClearance is how far a relation's caption sits off its line.
	relationLabelClearance styling.Length = 137160 // 0.15"
	// A caption's box grows with its text so it is not clipped.
	classRelationLabelWidth   styling.Length = 182880 // 0.2" of padding
	classRelationLabelPerRune styling.Length = 54864  // 0.06" per character
)

// relationLabelPosition places the caption of a class or ER relation clear of
// the boxes it runs between.
//
// Centring it on the connector's midpoint put it on top of a box whenever the
// two ends sat in the same row or column, which is the common case in a grid
// layout: the midpoint of "CUSTOMER ||--o{ ORDER" lands between them only if
// the gap happens to be wider than the caption.
//
// The caption is instead offset perpendicular to the run: above a horizontal
// connector, to the right of a vertical one, where nothing else is drawn.
func relationLabelPosition(
	startX, startY, endX, endY, labelWidth, labelHeight styling.Length,
) (styling.Length, styling.Length) {
	midX := (startX + endX) / 2
	midY := (startY + endY) / 2

	if absLength(endX-startX) >= absLength(endY-startY) {
		// Horizontal run: sit the caption above the line.
		return midX - labelWidth/2, midY - labelHeight - relationLabelClearance
	}
	// Vertical run: sit it to the right of the line.
	return midX + relationLabelClearance, midY - labelHeight/2
}

func absLength(value styling.Length) styling.Length {
	if value < 0 {
		return -value
	}
	return value
}

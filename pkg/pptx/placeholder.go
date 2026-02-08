package pptx

// PlaceholderType defines the type of placeholder.
type PlaceholderType string

const (
	PlaceholderTypeTitle      PlaceholderType = "title"
	PlaceholderTypeBody       PlaceholderType = "body"
	PlaceholderTypeCentrTitle PlaceholderType = "ctrTitle"
	PlaceholderTypeSubTitle   PlaceholderType = "subTitle"
	PlaceholderTypeDt         PlaceholderType = "dt"
	PlaceholderTypeSldNum     PlaceholderType = "sldNum"
	PlaceholderTypeFtr        PlaceholderType = "ftr"
	PlaceholderTypeHdr        PlaceholderType = "hdr"
	PlaceholderTypeObj        PlaceholderType = "obj"
	PlaceholderTypeChart      PlaceholderType = "chart"
	PlaceholderTypeTbl        PlaceholderType = "tbl"
	PlaceholderTypeClipArt    PlaceholderType = "clipArt"
	PlaceholderTypeDgm        PlaceholderType = "dgm"
	PlaceholderTypeMedia      PlaceholderType = "media"
	PlaceholderTypeSldImg     PlaceholderType = "sldImg"
	PlaceholderTypePic        PlaceholderType = "pic"
)

// Placeholder represents a placeholder shape on a slide layout or master.
// It can be populated with content (text, image, table, chart).
type Placeholder struct {
	Type  PlaceholderType
	Index int
	Name  string

	// underlying shape reference or geometry for layout inheritance
	X, Y, CX, CY int64
}

// InsertPicture inserts an image into the placeholder.
// The image is resized to fit the placeholder bounds while maintaining aspect ratio.
func (p *Placeholder) InsertPicture(imagePath string) Image {
	return Image{
		Path:        imagePath,
		X:           p.X,
		Y:           p.Y,
		CX:          p.CX,
		CY:          p.CY,
		Placeholder: p,
	}
}

// InsertPictureFromBytes inserts an image from bytes into the placeholder.
func (p *Placeholder) InsertPictureFromBytes(data []byte, format string) Image {
	return Image{
		Data:        data,
		Format:      format,
		X:           p.X,
		Y:           p.Y,
		CX:          p.CX,
		CY:          p.CY,
		Placeholder: p,
	}
}

// InsertText inserts text into the placeholder.
func (p *Placeholder) InsertText(text string) PlaceholderContent {
	return PlaceholderContent{
		Index: p.Index,
		Type:  string(p.Type),
		Text:  text,
	}
}

// InsertTable inserts a table into the placeholder.
func (p *Placeholder) InsertTable(table Table) PlaceholderContent {
	return PlaceholderContent{
		Index: p.Index,
		Type:  string(p.Type),
		Table: &table,
	}
}

// InsertPictureToSlide returns a PlaceholderContent for the image.
func (p *Placeholder) InsertPictureToSlide(image Image) PlaceholderContent {
	return PlaceholderContent{
		Index: p.Index,
		Type:  string(p.Type),
		Image: &image,
	}
}

package editorcommon

// Shape represents a simplified view of a slide shape for editing.
type Shape struct {
	ID   int
	Name string
	Type string
	Text string
	X, Y int
	W, H int

	PlaceholderIndex *int   `json:"PlaceholderIndex,omitempty"`
	PlaceholderType  string `json:"PlaceholderType,omitempty"`

	Fill   *ShapeFill
	Line   *ShapeLine
	Shadow *ShapeShadow

	Adjustments []ShapeAdjustment
}

// ShapeAdjustment represents one preset-geometry adjustment formula.
type ShapeAdjustment struct {
	Name    string
	Formula string
}

// ShapeSearchQuery filters shapes for editor-wide search.
type ShapeSearchQuery struct {
	NameContains  string
	TypeEquals    string
	TextContains  string
	CaseSensitive bool
}

// ShapeSearchResult identifies one matched shape and its slide index.
type ShapeSearchResult struct {
	SlideIndex int
	Shape      Shape
}

// Hyperlink holds properties for click or hover actions.
type Hyperlink struct {
	Address        *string `json:"address,omitempty"`
	Action         *string `json:"action,omitempty"`
	Tooltip        *string `json:"tooltip,omitempty"`
	TargetSlide    *int    `json:"target_slide,omitempty"`
	TargetJump     *string `json:"jump,omitempty"`
	Macro          *string `json:"macro,omitempty"`
	History        *bool   `json:"history,omitempty"`
	HighlightClick *bool   `json:"highlight_click,omitempty"`
	EndSound       *bool   `json:"end_sound,omitempty"`
}

// TextRun represents a single formatted text segment.
type TextRun struct {
	Text           string     `json:"text"`
	Bold           *bool      `json:"bold,omitempty"`
	Italic         *bool      `json:"italic,omitempty"`
	Underline      *string    `json:"underline,omitempty"`
	Strikethrough  *string    `json:"strikethrough,omitempty"`
	Subscript      *bool      `json:"subscript,omitempty"`
	Superscript    *bool      `json:"superscript,omitempty"`
	Color          *string    `json:"color,omitempty"`
	Highlight      *string    `json:"highlight,omitempty"`
	Font           *string    `json:"font,omitempty"`
	SizePt         *int       `json:"size_pt,omitempty"`
	Code           *bool      `json:"code,omitempty"`
	AllCaps        *bool      `json:"all_caps,omitempty"`
	SmallCaps      *bool      `json:"small_caps,omitempty"`
	OutlineColor   *string    `json:"outline_color,omitempty"`
	OutlineWidthPt *float64   `json:"outline_width_pt,omitempty"`
	Hyperlink      *Hyperlink `json:"hyperlink,omitempty"`
	HoverAction    *Hyperlink `json:"hover_action,omitempty"`
}

// ShapeProps defines optional properties when creating a shape.
type ShapeProps struct {
	Name string    `json:"name,omitempty"`
	Runs []TextRun `json:"runs,omitempty"`
}

// TextFrame defines formatting properties for the text body container within a shape.
type TextFrame struct {
	MarginTop     *int     `json:"margin_top,omitempty"`
	MarginBottom  *int     `json:"margin_bottom,omitempty"`
	MarginLeft    *int     `json:"margin_left,omitempty"`
	MarginRight   *int     `json:"margin_right,omitempty"`
	WordWrap      *bool    `json:"word_wrap,omitempty"`
	AutoFit       *bool    `json:"auto_fit,omitempty"`      // Deprecated: use auto_fit_type instead
	AutoFitType   *string  `json:"auto_fit_type,omitempty"` // "none", "normal", "shape"
	VerticalAlign *string  `json:"vertical_align,omitempty"`
	Orientation   *string  `json:"orientation,omitempty"`
	Columns       *int     `json:"columns,omitempty"`
	Rotation      *float64 `json:"rotation,omitempty"` // Degrees, converted to OOXML 1/60000 degree units.
}

// Paragraph defines paragraph-level formatting controls.
type Paragraph struct {
	Indent         *int    `json:"indent,omitempty"`
	Hanging        *int    `json:"hanging,omitempty"`
	TabStops       []int   `json:"tab_stops,omitempty"`
	Alignment      *string `json:"alignment,omitempty"`
	Level          *int    `json:"level,omitempty"`
	LineSpacingPct *int    `json:"line_spacing_pct,omitempty"`
	LineSpacingPts *int    `json:"line_spacing_pts,omitempty"`
	SpaceBeforePts *int    `json:"space_before_pts,omitempty"`
	SpaceAfterPts  *int    `json:"space_after_pts,omitempty"`
	BulletStyle    *string `json:"bullet_style,omitempty"`
	BulletChar     *string `json:"bullet_char,omitempty"`
	BulletColor    *string `json:"bullet_color,omitempty"`
	BulletSizePct  *int    `json:"bullet_size_pct,omitempty"`
}

// ShapeFill defines generic shape fill controls.
type ShapeFill struct {
	Solid        *string        `json:"solid,omitempty"`
	Transparency *float64       `json:"transparency,omitempty"`
	Gradient     *GradientFill  `json:"gradient,omitempty"`
	Pattern      *PatternedFill `json:"pattern,omitempty"`
	Background   *bool          `json:"background,omitempty"`
}

// ShapeLine defines generic shape line controls.
type ShapeLine struct {
	Color            *string `json:"color,omitempty"`
	WidthEmu         *int    `json:"width_emu,omitempty"`
	DashStyle        *string `json:"dash_style,omitempty"`
	StartArrow       *string `json:"start_arrow,omitempty"`
	StartArrowWidth  *string `json:"start_arrow_width,omitempty"`
	StartArrowLength *string `json:"start_arrow_length,omitempty"`
	EndArrow         *string `json:"end_arrow,omitempty"`
	EndArrowWidth    *string `json:"end_arrow_width,omitempty"`
	EndArrowLength   *string `json:"end_arrow_length,omitempty"`
}

// ShapeShadow defines generic shape shadow controls.
type ShapeShadow struct {
	Inherit     *bool    `json:"inherit,omitempty"`
	Color       *string  `json:"color,omitempty"`
	BlurEmu     *int     `json:"blur_emu,omitempty"`
	DistanceEmu *int     `json:"distance_emu,omitempty"`
	AngleDeg    *float64 `json:"angle_deg,omitempty"`
}

// ShapeGlow defines generic shape glow controls.
type ShapeGlow struct {
	Color     *string `json:"color,omitempty"`
	RadiusEmu *int    `json:"radius_emu,omitempty"`
}

// ShapeBlur defines generic shape blur controls.
type ShapeBlur struct {
	RadiusEmu *int `json:"radius_emu,omitempty"`
}

// ShapeSoftEdge defines generic shape soft-edge controls.
type ShapeSoftEdge struct {
	RadiusEmu *int `json:"radius_emu,omitempty"`
}

// ShapeReflection defines generic shape reflection controls.
type ShapeReflection struct {
	BlurEmu     *int `json:"blur_emu,omitempty"`
	DistanceEmu *int `json:"distance_emu,omitempty"`
}

type GradientStop struct {
	PositionPct *float64 `json:"position_pct,omitempty"`
	Color       string   `json:"color"`
}
type GradientFill struct {
	AngleDeg *float64       `json:"angle_deg,omitempty"`
	Stops    []GradientStop `json:"stops,omitempty"`
}
type PatternedFill struct {
	Preset  *string `json:"preset,omitempty"`
	FgColor *string `json:"fg_color,omitempty"`
	BgColor *string `json:"bg_color,omitempty"`
}

type ImageMetadata struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Format      string `json:"format"`
	ContentType string `json:"content_type,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

type ImageCrop struct {
	Left   float64 `json:"left,omitempty"`
	Right  float64 `json:"right,omitempty"`
	Top    float64 `json:"top,omitempty"`
	Bottom float64 `json:"bottom,omitempty"`
}

type ShapeUpdate struct {
	Text        *string          `json:"text,omitempty"`
	Runs        *[]TextRun       `json:"runs,omitempty"`
	TextFrame   *TextFrame       `json:"text_frame,omitempty"`
	Paragraph   *Paragraph       `json:"paragraph,omitempty"`
	Fill        *ShapeFill       `json:"fill,omitempty"`
	Line        *ShapeLine       `json:"line,omitempty"`
	Shadow      *ShapeShadow     `json:"shadow,omitempty"`
	Glow        *ShapeGlow       `json:"glow,omitempty"`
	Blur        *ShapeBlur       `json:"blur,omitempty"`
	SoftEdge    *ShapeSoftEdge   `json:"soft_edge,omitempty"`
	Reflection  *ShapeReflection `json:"reflection,omitempty"`
	ClickAction *Hyperlink       `json:"click_action,omitempty"`
	HoverAction *Hyperlink       `json:"hover_action,omitempty"`
	Crop        *ImageCrop       `json:"crop,omitempty"`
	Rotation    *float64         `json:"rotation,omitempty"`
	FlipH       *bool            `json:"flip_h,omitempty"`
	FlipV       *bool            `json:"flip_v,omitempty"`
	X           *int             `json:"x,omitempty"`
	Y           *int             `json:"y,omitempty"`
	W           *int             `json:"w,omitempty"`
	H           *int             `json:"h,omitempty"`
}

// SlideImageRef describes one image relationship on a slide.
type SlideImageRef struct {
	Index  int
	RelID  string
	Target string
}

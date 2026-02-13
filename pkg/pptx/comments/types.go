package comments

import "time"

// Author represents a person making comments in the presentation.
type Author struct {
	ID         int64
	Name       string
	Initials   string
	UserID     string // Optional, GUID or similar
	ProviderID string // Optional, identity provider
	LastIndex  int    // Metadata for tracking comment index
	ColorIndex int    // Visual color index (0-based)
}

// Comment represents a single comment on a slide.
type Comment struct {
	ID       int64
	AuthorID int64
	Text     string
	Date     time.Time
	X, Y     int64 // EMU coordinates
	Index    int   // Sequential index for this author
}

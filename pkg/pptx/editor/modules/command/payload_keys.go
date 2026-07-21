package command

// Payload keys shared across command request and response construction. They
// are part of the JSON command protocol, so they are named once here rather
// than repeated as literals at each site.
const (
	KeyText       = "text"
	KeyRuns       = "runs"
	KeyProperties = "properties"
	KeyIndex      = "index"

	KeyMimeType = "mime_type"
	KeyPath     = "path"
	KeyData     = "data"
)

// Media insert labels used to describe the primary/secondary payload of a
// media insert spec in validation errors.
const (
	labelVideo  = "video"
	labelPoster = "poster"
	labelIcon   = "icon"
)

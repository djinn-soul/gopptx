package structural

import "fmt"

// Severity represents the impact level of a diagnostic issue.
type Severity int

const (
	// SeverityInfo indicates a minor issue or informative note.
	SeverityInfo Severity = iota
	// SeverityWarning indicates a potential issue that might cause problems in some readers.
	SeverityWarning
	// SeverityError indicates a critical issue that makes the file invalid or unreadable.
	SeverityError
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// IssueCode identifies the specific type of diagnostic issue.
type IssueCode string

const (
	CodeMissingPart          IssueCode = "MISSING_PART"
	CodeInvalidXml           IssueCode = "INVALID_XML"
	CodeBrokenRelationship   IssueCode = "BROKEN_RELATIONSHIP"
	CodeMissingSlideRef      IssueCode = "MISSING_SLIDE_REF"
	CodeOrphanSlide          IssueCode = "ORPHAN_SLIDE"
	CodeInvalidContentType   IssueCode = "INVALID_CONTENT_TYPE"
	CodeCorruptedEntry       IssueCode = "CORRUPTED_ENTRY"
	CodeMissingNamespace     IssueCode = "MISSING_NAMESPACE"
	CodeEmptyRequiredElement IssueCode = "EMPTY_REQUIRED_ELEMENT"
	CodeModelValidationError IssueCode = "MODEL_VALIDATION_ERROR"
)

// Issue represents a single diagnostic finding.
type Issue struct {
	Code        IssueCode
	Severity    Severity
	Path        string // The file path within the PPTX package
	Description string
	Repairable  bool
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s (%s): %s (Repairable: %v)", i.Severity, i.Code, i.Path, i.Description, i.Repairable)
}

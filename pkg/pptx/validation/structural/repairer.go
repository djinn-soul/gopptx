package structural

import (
	"fmt"
)

// PartModifier defines the interface for modifying package parts.
type PartModifier interface {
	PartProvider
	Set(path string, data []byte)
	Delete(path string)
}

// Repairer provides methods for repairing detected diagnostic issues.
type Repairer struct {
	modifier PartModifier
}

// NewRepairer creates a new repairer using the given part modifier.
func NewRepairer(modifier PartModifier) *Repairer {
	return &Repairer{
		modifier: modifier,
	}
}

// RepairResult summarizes the outcome of a repair operation.
type RepairResult struct {
	IssuesRepaired   []Issue
	IssuesUnrepaired []Issue
}

// Repair attempts to fix the provided list of issues.
func (r *Repairer) Repair(issues []Issue) RepairResult {
	result := RepairResult{}
	for _, issue := range issues {
		if !issue.Repairable {
			result.IssuesUnrepaired = append(result.IssuesUnrepaired, issue)
			continue
		}

		if err := r.repairIssue(issue); err != nil {
			result.IssuesUnrepaired = append(result.IssuesUnrepaired, issue)
		} else {
			result.IssuesRepaired = append(result.IssuesRepaired, issue)
		}
	}
	return result
}

func (r *Repairer) repairIssue(issue Issue) error {
	switch issue.Code {
	case CodeMissingPart:
		return r.repairMissingPart(issue.Path)
	case CodeInvalidXML:
		return r.repairInvalidXML(issue.Path)
	case CodeBrokenRelationship:
		return r.repairBrokenRelationship(issue)
	case CodeOrphanSlide:
		return r.repairOrphanSlide(issue.Path)
	case CodeMissingSlideRef:
		return r.repairMissingSlideRef(issue)
	case CodeInvalidContentType:
		return r.repairInvalidContentType(issue.Path)
	case CodeMissingNamespace:
		return r.repairMissingNamespace(issue)
	case CodeEmptyRequiredElement:
		return r.repairEmptyRequiredElement(issue)
	default:
		return fmt.Errorf("unsupported repair: %s", issue.Code)
	}
}

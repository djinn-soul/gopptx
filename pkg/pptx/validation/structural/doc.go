// Package structural provides low-level PPTX package validation and repair flows.
//
// Validation strategy:
//   - Validate required package parts (`[Content_Types].xml`, `presentation.xml`, rels).
//   - Validate XML well-formedness for all XML and `.rels` parts.
//   - Run part-level namespace, relationship, and empty-element checks in parallel.
//   - Validate relationship targets and flag broken internal references.
//   - Validate presentation slide references and content-type registrations.
//   - Return structured issues with severity, repairability, and optional context.
//
// Repair strategy:
//   - Attempt deterministic, safe repairs only for known issue classes.
//   - Repair malformed XML (for example bare ampersands), broken relationships,
//     missing-or-invalid content-type registrations, and orphan-slide references.
//   - Preserve non-repairable issues in the output so callers can decide whether
//     to fail, warn, or continue.
//
// The repair APIs are designed for tooling and CLI integration. Use [Validator]
// followed by [Repairer] to implement validate/repair pipelines.
package structural

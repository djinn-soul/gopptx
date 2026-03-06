package command

import (
	"encoding/base64"
	"fmt"
)

type OptionalStringFieldFn func(map[string]any, string) string
type OptionalBoolFieldFn func(map[string]any, string) (bool, bool)
type AddValidationErrorFn func(code, message string)

type FindReplaceRequest struct {
	Find    string
	Replace string
}

func ParseFindReplaceRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
) (FindReplaceRequest, bool) {
	find, ok := parseStringField(payload, "find")
	if !ok {
		return FindReplaceRequest{}, false
	}
	replace, ok := parseStringField(payload, "replace")
	if !ok {
		return FindReplaceRequest{}, false
	}
	return FindReplaceRequest{
		Find:    find,
		Replace: replace,
	}, true
}

type AuthorAddRequest struct {
	Name     string
	Initials string
}

func ParseAuthorAddRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
) (AuthorAddRequest, bool) {
	name, ok := parseStringField(payload, "name")
	if !ok {
		return AuthorAddRequest{}, false
	}
	initials, ok := parseStringField(payload, "initials")
	if !ok {
		return AuthorAddRequest{}, false
	}
	return AuthorAddRequest{
		Name:     name,
		Initials: initials,
	}, true
}

type CommentAddRequest struct {
	SlideIndex int
	AuthorID   int64
	Text       string
	X          int64
	Y          int64
}

func ParseCommentAddRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseInt64Field ParseInt64FieldFn,
	parseStringField ParseStringFieldFn,
) (CommentAddRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return CommentAddRequest{}, false
	}
	authorID, ok := parseInt64Field(payload, "author_id")
	if !ok {
		return CommentAddRequest{}, false
	}
	text, ok := parseStringField(payload, "text")
	if !ok {
		return CommentAddRequest{}, false
	}
	x, ok := parseInt64Field(payload, "x")
	if !ok {
		return CommentAddRequest{}, false
	}
	y, ok := parseInt64Field(payload, "y")
	if !ok {
		return CommentAddRequest{}, false
	}
	return CommentAddRequest{
		SlideIndex: slideIndex,
		AuthorID:   authorID,
		Text:       text,
		X:          x,
		Y:          y,
	}, true
}

type CommentRemoveRequest struct {
	SlideIndex  int
	AuthorID    int64
	AuthorIndex int
}

func ParseCommentRemoveRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseInt64Field ParseInt64FieldFn,
	parseIntField ParseIntFieldFn,
) (CommentRemoveRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return CommentRemoveRequest{}, false
	}
	authorID, ok := parseInt64Field(payload, "author_id")
	if !ok {
		return CommentRemoveRequest{}, false
	}
	authorIndex, ok := parseIntField(payload, "author_index")
	if !ok {
		return CommentRemoveRequest{}, false
	}
	return CommentRemoveRequest{
		SlideIndex:  slideIndex,
		AuthorID:    authorID,
		AuthorIndex: authorIndex,
	}, true
}

type SetModifyPasswordRequest struct {
	Password string
}

func ParseSetModifyPasswordRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
) (SetModifyPasswordRequest, bool) {
	password, ok := parseStringField(payload, "password")
	if !ok {
		return SetModifyPasswordRequest{}, false
	}
	return SetModifyPasswordRequest{Password: password}, true
}

func ParseSetMarkAsFinalRequest(
	payload map[string]any,
	parseOptionalBool OptionalBoolFieldFn,
) (bool, bool) {
	return parseOptionalBool(payload, "final")
}

type CustomXMLAddRequest struct {
	Content     string
	RootElement string
	Namespace   string
	Properties  map[string]string
}

func ParseCustomXMLAddRequest(
	payload map[string]any,
	optionalString OptionalStringFieldFn,
	addErr AddValidationErrorFn,
	missingFieldCode string,
	invalidTypeCode string,
) CustomXMLAddRequest {
	content := optionalString(payload, "content")
	rootElement := optionalString(payload, "root_element")
	namespace := optionalString(payload, "namespace")
	properties := ParseCustomXMLProperties(payload, addErr, invalidTypeCode)

	if content == "" && rootElement == "" {
		addErr(missingFieldCode, "either content or root_element must be provided")
	}

	return CustomXMLAddRequest{
		Content:     content,
		RootElement: rootElement,
		Namespace:   namespace,
		Properties:  properties,
	}
}

func ParseCustomXMLProperties(
	payload map[string]any,
	addErr AddValidationErrorFn,
	invalidTypeCode string,
) map[string]string {
	rawProps, ok := payload["properties"]
	if !ok || rawProps == nil {
		return nil
	}
	propMap, ok := rawProps.(map[string]any)
	if !ok {
		addErr(invalidTypeCode, "properties must be an object with string values")
		return nil
	}
	props := make(map[string]string, len(propMap))
	for key, value := range propMap {
		s, ok := value.(string)
		if !ok {
			addErr(invalidTypeCode, fmt.Sprintf("property %q must be a string", key))
			continue
		}
		props[key] = s
	}
	return props
}

func DecodeRequiredBase64Field(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
	field string,
	invalidDataErr string,
) ([]byte, bool, error) {
	value, ok := parseStringField(payload, field)
	if !ok {
		return nil, false, nil
	}
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, true, fmt.Errorf("%s", invalidDataErr)
	}
	return data, true, nil
}

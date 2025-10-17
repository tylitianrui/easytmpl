package easytmpl

import (
	"bytes"
	"errors"
)

var (
	// TagContainSpaceError indicates that a tag contains spaces, which is not allowed.
	TagContainSpaceError = errors.New("tag cannot contain spaces")

	// TagEmptyError indicates that a tag is an empty string, which is not allowed.
	TagEmptyError = errors.New("tag cannot be an empty string")
)

// DefaultTagPair default tag pair `{{` and `}}`.
var DefaultTagPair = &TagPair{
	start: []byte{'{', '{'},
	end:   []byte{'}', '}'},
}

// TagPair defines a pair of tags to denote the start and end of a placeholder in the template.
// For example, in the template "{{name}}", the start tag is "{{" and the end tag is "}}".
type TagPair struct {
	start []byte
	end   []byte
}

// NewTagPair creates a new TagPair with the specified start and end tags.
func NewTagPair(start, end string) (*TagPair, error) {
	s := s2b(start)
	if err := checkTagInvalid(s); err != nil {
		return nil, err
	}

	e := s2b(end)
	if err := checkTagInvalid(e); err != nil {
		return nil, err
	}

	tag := &TagPair{
		start: s,
		end:   e,
	}
	return tag, nil
}

// checkTagInvalid checks if a tag is valid (not empty and does not contain spaces).
// invalid tag will return an error.
func checkTagInvalid(tag []byte) error {
	if len(tag) == 0 {
		return TagEmptyError
	}
	if bytes.Contains(tag, []byte{' '}) {
		return TagContainSpaceError
	}
	return nil
}

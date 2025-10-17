package easytmpl

import (
	"errors"
	"math"
)

// OptionHandler defines a function type for configuring Template options.
type OptionHandler func(*Template) error

// WithTagPair sets a custom tag pair for the template.
func WithTagPair(start, end string) OptionHandler {
	return func(t *Template) error {
		tag, err := NewTagPair(start, end)
		if err != nil {
			return err
		}
		t.pairs = tag
		return nil
	}
}

// WithPreAllocateMemory sets the initial capacity for the internal buffer used during template rendering.
func WithPreAllocateMemory(n int) OptionHandler {
	return func(t *Template) error {
		if n <= 0 || n > math.MaxInt {
			return errors.New("invalid pre-allocated memory size")
		}
		t.capacity = n
		return nil
	}
}

// WithAutoFill sets a default value to automatically fill in for any missing parameters during template rendering.
// it only works in non-strict mode.
func WithAutoFill(s string) OptionHandler {
	return func(t *Template) error {
		b := s2b(s)
		t.autoFill = &b
		return nil
	}
}

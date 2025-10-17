package easytmpl

import (
	"bytes"
	"errors"
	"io"
	"math"
)

var (
	// TemplateContentEmptyError indicates that the template content is empty.
	TemplateContentEmptyError = errors.New("template content is empty")

	// TemplateExecMissingParameterError indicates that a required parameter is missing during template rendering process .
	TemplateExecMissingParameterError = errors.New("missing parameter")
)

// Template implements a template engine, which supports custom tags(aka placeholders) and parameters rendering.
type Template struct {
	content            []byte
	contentIntervalIdx [][2]int
	args               [][]byte
	pairs              *TagPair
	capacity           int
	autoFill           *[]byte
}

// NewTemplate creates a new Template instance with the provided template string and optional configurations.
// If no tag pair is specified, the default tag pair `{{` and `}}` will be used.
// It returns an error if the template content is empty or consists solely of whitespace.
func NewTemplate(tpl string, opts ...OptionHandler) (*Template, error) {

	if len(tpl) == 0 {
		return nil, TemplateContentEmptyError
	}

	content := s2b(tpl)

	var isAllBlank = true
	for _, c := range content {
		if c != ' ' {
			isAllBlank = false
			break
		}
	}
	if isAllBlank {
		return nil, TemplateContentEmptyError
	}

	template := &Template{
		content: content,
	}
	for _, opt := range opts {
		if err := opt(template); err != nil {
			return nil, err
		}
	}

	if template.pairs == nil {
		template.pairs = DefaultTagPair
	}
	template.parse()
	return template, nil
}

// parse parses the template content to identify placeholders and their positions based on the defined tag pairs.
// It populates the args slice with the identified placeholders and
// the contentIntervalIdx slice with the intervals of static content.
func (t *Template) parse() {
	slen := len(t.pairs.start)
	elen := len(t.pairs.end)
	en := len(t.content) - elen
	sn := len(t.content) - slen - elen
	if en <= 0 {
		t.args = nil
		t.contentIntervalIdx = nil
		return
	}

	var argStartIdx, argEndIdx = 0 - elen - 1, 0 - elen
	var lastStartIdx, lastEndIdx = argStartIdx, argEndIdx
	var j int

	for i := 0; i <= en; i++ {

		if i < sn && bytes.Equal(t.content[i:i+slen], t.pairs.start) {
			j = argStartIdx
			argStartIdx = i
			continue
		}

		if bytes.Equal(t.content[i:i+elen], t.pairs.end) {
			if lastEndIdx > argStartIdx {
				continue
			}
			if i > argStartIdx && argStartIdx > lastStartIdx {
				if !IsBlank(t.content[argStartIdx+slen : i]) {
					t.contentIntervalIdx = append(t.contentIntervalIdx, [2]int{argEndIdx + elen, argStartIdx})
					t.args = append(t.args, t.content[argStartIdx+slen:i])
					lastStartIdx = argStartIdx
					argEndIdx = i
				} else if j >= 0 {
					t.args = append(t.args, t.content[j+slen:i])
					t.contentIntervalIdx = append(t.contentIntervalIdx, [2]int{argEndIdx + elen, j})

					lastStartIdx = j
					argEndIdx = i

				}

			}
			lastEndIdx = i
		}
	}
	t.contentIntervalIdx = append(t.contentIntervalIdx, [2]int{argEndIdx + elen, math.MaxInt})

}

// ExecString renders the template with the provided arguments.
// If strict is true, it returns an error if any placeholder in the template
// does not have a corresponding entry in args.
// If strict is false, placeholders without corresponding entries in args will remain unchanged in the output.
func (t *Template) ExecString(args map[string]string, strict bool) (string, error) {
	if strict {
		for _, a := range t.args {
			if _, ok := args[string(a)]; !ok {
				return "", TemplateExecMissingParameterError
			}
		}
	}
	var bb bytes.Buffer
	if len(t.content) > t.capacity {
		bb.Grow(len(t.content) * 2)
	} else {
		bb.Grow(t.capacity)
	}

	err := t.exec(&bb, func(w io.Writer, key string) (int, error) {
		if v, ok := args[key]; ok {
			return w.Write(s2b(v))
		} else if t.autoFill != nil {
			return w.Write(*t.autoFill)
		} else {
			w.Write(t.pairs.start)
			w.Write(s2b(key))
			w.Write(t.pairs.end)
			return len(t.pairs.start) + len(t.pairs.end) + len(key), nil
		}
	})

	return bb.String(), err
}

// exec is a helper function that executes the template rendering process.
func (t *Template) exec(b io.Writer, f func(w io.Writer, key string) (int, error)) error {

	for i := 0; i < len(t.contentIntervalIdx)-1; i++ {
		c := t.content[t.contentIntervalIdx[i][0]:t.contentIntervalIdx[i][1]]
		b.Write(c)
		_, err := f(b, b2s(t.args[i]))
		if err != nil {
			return err
		}
	}

	c := t.content[t.contentIntervalIdx[len(t.contentIntervalIdx)-1][0]:]
	b.Write(c)

	return nil
}

// ExecuteFunc renders the template using a custom function to handle each placeholder.
// The function f is called for each placeholder with the writer and the placeholder key.
// It returns the rendered string or an error if any occurs during the rendering process.
func (t *Template) ExecuteFunc(w io.Writer, f func(w io.Writer, key string) (int, error)) error {
	return t.exec(w, f)
}

package timing

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	easytmpl "github.com/tylitianrui/easytmpl"
	fasttemplate "github.com/valyala/fasttemplate"
)

var args = map[string]string{
	"domain":   "example",
	"uid":      "12345",
	"name":     "tyltr",
	"from":     "NY",
	"to":       "LA",
	"birthday": "2024-01-01",
	"language": "en",
	"ref":      "go.template",
	"t":        "169611111199900000222",
}

var (

	// TemplateWith10PlaceholderTemplate with 10 placeholders.
	TemplateWith10Placeholder = "https://{{domain}}.com?uid={{uid}}&name={{name}}&from={{from}}&to={{to}}&birthday={{birthday}}&lang={{language}}&ref={{ref}}&t={{t}}"

	// ExpectedResultTemplateWith10Placeholderexpected result of template with 10 placeholders after rendering.
	ExpectedResultTemplateWith10Placeholder = "https://example.com?uid=12345&name=tyltr&from=NY&to=LA&birthday=2024-01-01&lang=en&ref=go.template&t=169611111199900000222"

	// ExpectedResultTemplateWith10PlaceholderBytes
	ExpectedResultTemplateWith10PlaceholderBytes = []byte(ExpectedResultTemplateWith10Placeholder)

	// Template with 20 placeholders
	TemplateWith20Placeholder = TemplateWith10Placeholder + TemplateWith10Placeholder

	// result of template with 20 placeholders after rendering
	ExpectedResultTemplateWith20Placeholder = ExpectedResultTemplateWith10Placeholder + ExpectedResultTemplateWith10Placeholder

	// ExpectedResultTemplateWith20PlaceholderBytes
	ExpectedResultTemplateWith20PlaceholderBytes = []byte(ExpectedResultTemplateWith20Placeholder)

	// Template with 30 placeholders
	TemplateWith30Placeholder = TemplateWith10Placeholder + TemplateWith10Placeholder + TemplateWith10Placeholder

	// result of template with 30 placeholders after rendering
	ExpectedResultTemplateWith30Placeholder = ExpectedResultTemplateWith10Placeholder + ExpectedResultTemplateWith10Placeholder + ExpectedResultTemplateWith10Placeholder

	// ExpectedResultTemplateWith30PlaceholderBytes
	ExpectedResultTemplateWith30PlaceholderBytes = []byte(ExpectedResultTemplateWith30Placeholder)
)

func rendingFunc(w io.Writer, key string) (int, error) {
	if t, ok := args[key]; ok {
		return w.Write([]byte(t))
	}
	return 0, nil
}

func Benchmark_GoTextTemplate_ExecuteWith10Placeholder(b *testing.B) {
	tpl1 := strings.Replace(TemplateWith10Placeholder, "{{", "{{.", -1)
	t, err := template.New("go text template").Parse(tpl1)
	if err != nil {
		b.Fatalf("Error when parsing template: %s", err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if err := t.Execute(&w, args); err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			x := w.String()
			if x != ExpectedResultTemplateWith10Placeholder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith10Placeholder)
			}
			w.Reset()
		}
	})

}

func Benchmark_FastTemplate_ExecuteStringWith10Placeholder(b *testing.B) {

	t, err := fasttemplate.NewTemplate(TemplateWith10Placeholder, "{{", "}}")
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}
	mm := make(map[string]interface{}, len(args))
	for k, v := range args {
		mm[k] = v
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			x := t.ExecuteString(mm)
			if x != ExpectedResultTemplateWith10Placeholder {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith10Placeholder)
			}
		}
	})
}

func Benchmark_EasyTmpl_ExecuteStringWith10Placeholder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith10Placeholder)
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x, err := t.ExecString(args, true)
			if err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			if x != ExpectedResultTemplateWith10Placeholder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith10Placeholder)
			}
		}
	})
}

func Benchmark_FastTemplate_ExecuteFuncWith10Placeholder(b *testing.B) {
	t, err := fasttemplate.NewTemplate(TemplateWith10Placeholder, "{{", "}}")
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if _, err := t.ExecuteFunc(&w, rendingFunc); err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			x := w.Bytes()
			if !bytes.Equal(x, ExpectedResultTemplateWith10PlaceholderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith10PlaceholderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_EasyTmpl_ExecuteFuncWith10Placeholder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith10Placeholder)
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	mm := make(map[string]interface{})
	for k, v := range args {
		mm[k] = v
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			t.ExecuteFunc(&w, rendingFunc)
			x := w.Bytes()
			if !bytes.Equal(x, ExpectedResultTemplateWith10PlaceholderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith10PlaceholderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_GoTextTemplate_ExecuteWith20Placeholder(b *testing.B) {
	tpl := strings.Replace(TemplateWith20Placeholder, "{{", "{{.", -1)
	t, err := template.New("go text template").Parse(tpl)
	if err != nil {
		b.Fatalf("Error when parsing template: %s", err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if err := t.Execute(&w, args); err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			x := w.String()
			if x != ExpectedResultTemplateWith20Placeholder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith20Placeholder)
			}
			w.Reset()
		}
	})

}

func Benchmark_FastTemplate_ExecuteStringWith20Placeholder(b *testing.B) {

	t, err := fasttemplate.NewTemplate(TemplateWith20Placeholder, "{{", "}}")
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}
	mm := make(map[string]interface{}, len(args))
	for k, v := range args {
		mm[k] = v
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			x := t.ExecuteString(mm)
			if x != ExpectedResultTemplateWith20Placeholder {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith20Placeholder)
			}
		}
	})
}

func Benchmark_EasyTmpl_ExecuteStringWith20Placeholder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith20Placeholder)
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x, err := t.ExecString(args, true)
			if err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			if x != ExpectedResultTemplateWith20Placeholder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith20Placeholder)
			}
		}
	})
}

func Benchmark_FastTemplate_ExecuteFuncWith20Placeholder(b *testing.B) {
	t, err := fasttemplate.NewTemplate(TemplateWith20Placeholder, "{{", "}}")
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if _, err := t.ExecuteFunc(&w, rendingFunc); err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			x := w.Bytes()
			if !bytes.Equal(x, ExpectedResultTemplateWith20PlaceholderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith20PlaceholderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_EasyTmpl_ExecuteFuncWith20Placeholder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith20Placeholder)
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	mm := make(map[string]interface{})
	for k, v := range args {
		mm[k] = v
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			t.ExecuteFunc(&w, rendingFunc)
			x := w.Bytes()
			if !bytes.Equal(x, ExpectedResultTemplateWith20PlaceholderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith20PlaceholderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_GoTextTemplate_ExecuteWith30Placeholder(b *testing.B) {
	tpl := strings.Replace(TemplateWith30Placeholder, "{{", "{{.", -1)
	t, err := template.New("go text template").Parse(tpl)
	if err != nil {
		b.Fatalf("Error when parsing template: %s", err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if err := t.Execute(&w, args); err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			x := w.String()
			if x != ExpectedResultTemplateWith30Placeholder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith30Placeholder)
			}
			w.Reset()
		}
	})

}

func Benchmark_FastTemplate_ExecuteStringWith30Placeholder(b *testing.B) {

	t, err := fasttemplate.NewTemplate(TemplateWith30Placeholder, "{{", "}}")
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}
	mm := make(map[string]interface{}, len(args))
	for k, v := range args {
		mm[k] = v
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			x := t.ExecuteString(mm)
			if x != ExpectedResultTemplateWith30Placeholder {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith30Placeholder)
			}
		}
	})
}

func Benchmark_EasyTmpl_ExecuteStringWith30Placeholder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith30Placeholder)
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x, err := t.ExecString(args, true)
			if err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			if x != ExpectedResultTemplateWith30Placeholder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith30Placeholder)
			}
		}
	})
}

func Benchmark_FastTemplate_ExecuteFuncWith30Placeholder(b *testing.B) {
	t, err := fasttemplate.NewTemplate(TemplateWith30Placeholder, "{{", "}}")
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if _, err := t.ExecuteFunc(&w, rendingFunc); err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			x := w.Bytes()
			if !bytes.Equal(x, ExpectedResultTemplateWith30PlaceholderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith30PlaceholderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_EasyTmpl_ExecuteFuncWith30Placeholder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith30Placeholder)
	if err != nil {
		b.Fatalf("error in template: %s", err)
	}

	mm := make(map[string]interface{})
	for k, v := range args {
		mm[k] = v
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			t.ExecuteFunc(&w, rendingFunc)
			x := w.Bytes()
			if !bytes.Equal(x, ExpectedResultTemplateWith30PlaceholderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith30PlaceholderBytes)
			}
			w.Reset()
		}
	})
}

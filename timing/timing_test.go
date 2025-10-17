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

	// TemplateWith10SpaceHolder Template with 10 placeholders.
	TemplateWith10SpaceHolder = "https://{{domain}}.com?uid={{uid}}&name={{name}}&from={{from}}&to={{to}}&birthday={{birthday}}&lang={{language}}&ref={{ref}}&t={{t}}"

	// ExpectedResultTemplateWith10SpaceHolder expected result of template with 10 placeholders after rendering.
	ExpectedResultTemplateWith10SpaceHolder = "https://example.com?uid=12345&name=tyltr&from=NY&to=LA&birthday=2024-01-01&lang=en&ref=go.template&t=169611111199900000222"

	// ExpectedResultTemplateWith10SpaceHolderBytes
	ExpectedResultTemplateWith10SpaceHolderBytes = []byte(ExpectedResultTemplateWith10SpaceHolder)

	// Template with 20 placeholders
	TemplateWith20SpaceHolder = TemplateWith10SpaceHolder + TemplateWith10SpaceHolder

	// result of template with 20 placeholders after rendering
	ExpectedResultTemplateWith20SpaceHolder = ExpectedResultTemplateWith10SpaceHolder + ExpectedResultTemplateWith10SpaceHolder

	// ExpectedResultTemplateWith20SpaceHolderBytes
	ExpectedResultTemplateWith20SpaceHolderBytes = []byte(ExpectedResultTemplateWith20SpaceHolder)

	// Template with 30 placeholders
	TemplateWith30SpaceHolder = TemplateWith10SpaceHolder + TemplateWith10SpaceHolder + TemplateWith10SpaceHolder

	// result of template with 30 placeholders after rendering
	ExpectedResultTemplateWith30SpaceHolder = ExpectedResultTemplateWith10SpaceHolder + ExpectedResultTemplateWith10SpaceHolder + ExpectedResultTemplateWith10SpaceHolder

	// ExpectedResultTemplateWith30SpaceHolderBytes
	ExpectedResultTemplateWith30SpaceHolderBytes = []byte(ExpectedResultTemplateWith30SpaceHolder)
)

func rendingFunc(w io.Writer, key string) (int, error) {
	if t, ok := args[key]; ok {
		return w.Write([]byte(t))
	}
	return 0, nil
}

func Benchmark_GoTextTemplate_ExecuteWith10SpaceHolder(b *testing.B) {
	tpl1 := strings.Replace(TemplateWith10SpaceHolder, "{{", "{{.", -1)
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
			if x != ExpectedResultTemplateWith10SpaceHolder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith10SpaceHolder)
			}
			w.Reset()
		}
	})

}

func Benchmark_FastTemplate_ExecuteStringWith10SpaceHolder(b *testing.B) {

	t, err := fasttemplate.NewTemplate(TemplateWith10SpaceHolder, "{{", "}}")
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
			if x != ExpectedResultTemplateWith10SpaceHolder {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith10SpaceHolder)
			}
		}
	})
}

func Benchmark_EasyTmpl_ExecuteStringWith10SpaceHolder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith10SpaceHolder)
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
			if x != ExpectedResultTemplateWith10SpaceHolder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith10SpaceHolder)
			}
		}
	})
}

func Benchmark_FastTemplate_ExecuteFuncWith10SpaceHolder(b *testing.B) {
	t, err := fasttemplate.NewTemplate(TemplateWith10SpaceHolder, "{{", "}}")
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
			if !bytes.Equal(x, ExpectedResultTemplateWith10SpaceHolderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith10SpaceHolderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_EasyTmpl_ExecuteFuncWith10SpaceHolder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith10SpaceHolder)
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
			if !bytes.Equal(x, ExpectedResultTemplateWith10SpaceHolderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith10SpaceHolderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_GoTextTemplate_ExecuteWith20SpaceHolder(b *testing.B) {
	tpl := strings.Replace(TemplateWith20SpaceHolder, "{{", "{{.", -1)
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
			if x != ExpectedResultTemplateWith20SpaceHolder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith20SpaceHolder)
			}
			w.Reset()
		}
	})

}

func Benchmark_FastTemplate_ExecuteStringWith20SpaceHolder(b *testing.B) {

	t, err := fasttemplate.NewTemplate(TemplateWith20SpaceHolder, "{{", "}}")
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
			if x != ExpectedResultTemplateWith20SpaceHolder {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith20SpaceHolder)
			}
		}
	})
}

func Benchmark_EasyTmpl_ExecuteStringWith20SpaceHolder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith20SpaceHolder)
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
			if x != ExpectedResultTemplateWith20SpaceHolder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith20SpaceHolder)
			}
		}
	})
}

func Benchmark_FastTemplate_ExecuteFuncWith20SpaceHolder(b *testing.B) {
	t, err := fasttemplate.NewTemplate(TemplateWith20SpaceHolder, "{{", "}}")
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
			if !bytes.Equal(x, ExpectedResultTemplateWith20SpaceHolderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith20SpaceHolderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_EasyTmpl_ExecuteFuncWith20SpaceHolder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith20SpaceHolder)
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
			if !bytes.Equal(x, ExpectedResultTemplateWith20SpaceHolderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith20SpaceHolderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_GoTextTemplate_ExecuteWith30SpaceHolder(b *testing.B) {
	tpl := strings.Replace(TemplateWith30SpaceHolder, "{{", "{{.", -1)
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
			if x != ExpectedResultTemplateWith30SpaceHolder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith30SpaceHolder)
			}
			w.Reset()
		}
	})

}

func Benchmark_FastTemplate_ExecuteStringWith30SpaceHolder(b *testing.B) {

	t, err := fasttemplate.NewTemplate(TemplateWith30SpaceHolder, "{{", "}}")
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
			if x != ExpectedResultTemplateWith30SpaceHolder {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith30SpaceHolder)
			}
		}
	})
}

func Benchmark_EasyTmpl_ExecuteStringWith30SpaceHolder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith30SpaceHolder)
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
			if x != ExpectedResultTemplateWith30SpaceHolder {
				b.Fatalf("unexpected result\n%s\nExpected\n%s\n", x, ExpectedResultTemplateWith30SpaceHolder)
			}
		}
	})
}

func Benchmark_FastTemplate_ExecuteFuncWith30SpaceHolder(b *testing.B) {
	t, err := fasttemplate.NewTemplate(TemplateWith30SpaceHolder, "{{", "}}")
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
			if !bytes.Equal(x, ExpectedResultTemplateWith30SpaceHolderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith30SpaceHolderBytes)
			}
			w.Reset()
		}
	})
}

func Benchmark_EasyTmpl_ExecuteFuncWith30SpaceHolder(b *testing.B) {
	t, err := easytmpl.NewTemplate(TemplateWith30SpaceHolder)
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
			if !bytes.Equal(x, ExpectedResultTemplateWith30SpaceHolderBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, ExpectedResultTemplateWith30SpaceHolderBytes)
			}
			w.Reset()
		}
	})
}

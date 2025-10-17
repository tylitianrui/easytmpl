package easytmpl

import (
	"math"
	"reflect"
	"testing"
)

func TestTemplate_parse(t *testing.T) {
	tests := []struct {
		name     string
		template *Template
		want     *Template
	}{
		{
			name: "case:{{e}}",
			template: &Template{
				content: []byte("{{e}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{e}}"),
				contentIntervalIdx: [][2]int{{0, 0}, {5, math.MaxInt}},
				args:               [][]byte{{'e'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},
		{
			name: "case:{{}}",
			template: &Template{
				content: []byte("{{}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{}}"),
				contentIntervalIdx: [][2]int{{0, math.MaxInt}},
				args:               nil,
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},
		{
			name: "case:{{{e}}",
			template: &Template{
				content: []byte("{{{e}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{{e}}"),
				contentIntervalIdx: [][2]int{{0, 1}, {6, math.MaxInt}},
				args:               [][]byte{{'e'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},

		{
			name: "case:{{{e}}}}}",
			template: &Template{
				content: []byte("{{{e}}}}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{{e}}}}}"),
				contentIntervalIdx: [][2]int{{0, 1}, {6, math.MaxInt}},
				args:               [][]byte{{'e'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},

		{
			name: "case:{{{e}}{}}}}",
			template: &Template{
				content: []byte("{{{e}}{}}}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{{e}}{}}}}"),
				contentIntervalIdx: [][2]int{{0, 1}, {6, math.MaxInt}},
				args:               [][]byte{{'e'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},

		{
			name: "case:{{{e}}{{a1}}}}}",
			template: &Template{
				content: []byte("{{{e}}{{a1}}}}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{{e}}{{a1}}}}}"),
				contentIntervalIdx: [][2]int{{0, 1}, {6, 6}, {12, math.MaxInt}},
				args:               [][]byte{{'e'}, {'a', '1'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},

		{
			name: "case:}}{{a{{}}b}}{{c}}{{",
			template: &Template{
				content: []byte("}}{{a{{}}b}}{{c}}{{"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("}}{{a{{}}b}}{{c}}{{"),
				contentIntervalIdx: [][2]int{{0, 2}, {9, 12}, {17, math.MaxInt}},
				args:               [][]byte{{'a', '{', '{'}, {'c'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},
		{
			name: "case:{{{name{{}}tyltr}}",
			template: &Template{
				content: []byte("{{{name{{}}tyltr}}"),

				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
			want: &Template{
				content:            []byte("{{{name{{}}tyltr}}"),
				contentIntervalIdx: [][2]int{{0, 1}, {11, math.MaxInt}},
				args:               [][]byte{{'n', 'a', 'm', 'e', '{', '{'}},
				pairs: &TagPair{
					start: []byte{'{', '{'},
					end:   []byte{'}', '}'},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.template.parse()
			if !reflect.DeepEqual(tt.template, tt.want) {
				t.Fatalf("parse fail expect:%v  actual: %v", tt.want, tt.template)
			}

		})
	}
}

func TestTemplate_Exec(t *testing.T) {

	t.Run("case:default tag pair `{{` & `}}` with strict mode", func(t *testing.T) {
		txt := "i am {{name}}, {{age}} year old, from {{country}}"
		template, err := NewTemplate(txt)
		if err != nil {
			t.Errorf("error %v", err)
		}
		args := map[string]string{
			"name":    "tyltr",
			"age":     "18",
			"country": "china",
		}
		got, err := template.ExecString(args, true)
		if err != nil {
			t.Errorf("error %v", err)
		}
		want := "i am tyltr, 18 year old, from china"
		if got != want {
			t.Errorf("error %v", err)
		}
	})

	t.Run("case:user-defined tag pair `[[` & `]]` with strict mode", func(t *testing.T) {
		txt := "i am [[name]], [[age]] year old, from [[country]]"
		template, err := NewTemplate(txt, WithTagPair("[[", "]]"))
		if err != nil {
			t.Errorf("error %v", err)
		}
		args := map[string]string{
			"name":    "tyltr",
			"age":     "18",
			"country": "china",
		}
		got, err := template.ExecString(args, true)
		if err != nil {
			t.Errorf("error %v", err)
		}
		want := "i am tyltr, 18 year old, from china"
		if got != want {
			t.Errorf("error %v", err)
		}
	})
	t.Run("case:golang  official  tag `{{.xxxx}}` with strict mode ", func(t *testing.T) {
		txt := "i am {{.name}}, {{.age}} year old, from {{.country}}"
		template, err := NewTemplate(txt, WithTagPair("{{.", "}}"))
		if err != nil {
			t.Errorf("error %v", err)
		}
		args := map[string]string{
			"name":    "tyltr",
			"age":     "18",
			"country": "china",
		}
		got, err := template.ExecString(args, true)
		if err != nil {
			t.Errorf("error %v", err)
		}
		want := "i am tyltr, 18 year old, from china"
		if got != want {
			t.Errorf("error %v", err)
		}
	})

	t.Run("case:golang  official  tag `{{.xxxx}}` without strict mode", func(t *testing.T) {
		txt := "i am {{.name}}, {{.age}} year old, from {{.country}}"
		template, err := NewTemplate(txt, WithTagPair("{{.", "}}"))
		if err != nil {
			t.Errorf("error %v", err)
		}
		args := map[string]string{
			"age":     "18",
			"country": "china",
		}
		got, err := template.ExecString(args, false)
		if err != nil {
			t.Errorf("error %v", err)
		}
		want := "i am {{.name}}, 18 year old, from china"
		if got != want {
			t.Errorf("error %v", err)
		}
	})
}

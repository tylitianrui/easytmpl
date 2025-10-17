package easytmpl

import (
	"fmt"
)

func ExampleExecString_NonStrictMode() {
	tpl := "https://{{demain}}.com?name={{name}}&age={{age}}&birth={{birth}}"

	// Create a new template instance with default tag pair `{{ }}` and pre-allocated memory of 1024 bytes.
	t, err := NewTemplate(tpl,
		WithPreAllocateMemory(1024), // Pre-allocate memory for better performance ,1024 bytes
	)

	if err != nil {
		panic(err)
	}

	// Substitution map.
	// "birth" tag is missing in the map
	args := map[string]string{
		"demain": "user.google",
		"bar":    "foobar",
		"name":   "tyltr",
		"age":    "18",
		// "birth" is missing.
	}

	// Non-strict mode, placeholders without corresponding entries in args will remain unchanged in the output.
	s, err := t.ExecString(args, false)
	fmt.Println("err:", err)
	fmt.Println("template:", s)

	// Output:
	// err: <nil>
	// template: https://user.google.com?name=tyltr&age=18&birth={{birth}}
}

func ExampleExecString_NonStrictModeAndAutoFill() {
	tpl := "https://[[demain]].com?name=[[name]]&age=[[age]]&birth=[[birth]]"
	t, err := NewTemplate(tpl,
		WithTagPair("[[", "]]"),     // set custom tag pair `[[` & `]]`
		WithPreAllocateMemory(1024), // Pre-allocate memory for better performance ,1024 bytes
		WithAutoFill(""),            // Auto fill missing parameters with empty string
	)
	if err != nil {
		panic(err)
	}

	// Substitution map.
	// "birth" tag is missing in the map
	args := map[string]string{
		"demain": "user.google",
		"bar":    "foobar",
		"name":   "tyltr",
		"age":    "18",
		// "birth" is missing.
	}

	// Non-strict mode, placeholders without corresponding entries in args will remain unchanged in the output.
	s, err := t.ExecString(args, false)
	fmt.Println("err:", err)
	fmt.Println("template:", s)

	// Output:
	// err: <nil>
	// template: https://user.google.com?name=tyltr&age=18&birth=
}

func ExampleExecString_StrictMode() {
	tpl := "https://{{demain}}.com?name={{name}}&age={{age}}&birth={{birth}}"
	t, err := NewTemplate(tpl,
		WithTagPair("{{", "}}"),
		WithPreAllocateMemory(1024),
	)
	if err != nil {
		panic(err)
	}

	// Substitution map.
	// "birth" tag is missing in the map
	args := map[string]string{
		"demain": "user.google",
		"bar":    "foobar",
		"name":   "tyltr",
		"age":    "18",
		// "birth" is missing.
	}

	// strict mode, it returns an error if any placeholder in the template does not have a corresponding entry in args.
	s, err := t.ExecString(args, true)
	fmt.Println("err:", err)
	fmt.Println("template:", s)

	// Output:
	// err: missing parameter
	// template:
}

# easytmpl

[![Go Reference](https://pkg.go.dev/badge/github.com/tylitianrui/easytmpl)](https://pkg.go.dev/github.com/tylitianrui/easytmpl) [![Go Report](https://goreportcard.com/badge/github.com/tylitianrui/easytmpl)](https://goreportcard.com/report/github.com/tylitianrui/easytmpl)

[English](./README.md) | **[中文文档](./README-CN.md)**

简单、高效的`go` 语言模版渲染引擎


## BenchMark 
以下是 **[easytmpl](https://github.com/tylitianrui/easytmpl)** 与 [fasttemplate](https://github.com/valyala/fasttemplate)`text/template` 性能比较的基准测试结果


### benchmark results of  10 space holders 
```
  timing git:(master) ✗ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/tylitianrui/easytmpl/timing
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
Benchmark_GoTextTemplate_ExecuteWith10PlaceHolder-4               497427              2177 ns/op             464 B/op         20 allocs/op
Benchmark_FastTemplate_ExecuteStringWith10PlaceHolder-4          1920590               749.1 ns/op           443 B/op         12 allocs/op
Benchmark_easytmpl_ExecuteStringWith10PlaceHolder-4            3118250               554.4 ns/op           464 B/op          3 allocs/op
Benchmark_FastTemplate_ExecuteFuncWith10PlaceHolder-4            2383504               486.6 ns/op           104 B/op          9 allocs/op
Benchmark_easytmpl_ExecuteFuncWith10PlaceHolder-4              2812179               362.3 ns/op           104 B/op          9 allocs/op
PASS
ok      github.com/tylitianrui/easytmpl/timing        9.356s

```

### benchmark results of  20 space holders 

```
➜  timing git:(master) ✗ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/tylitianrui/easytmpl/timing
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
Benchmark_GoTextTemplate_ExecuteWith20PlaceHolder-4               329619              3404 ns/op             880 B/op         38 allocs/op
Benchmark_FastTemplate_ExecuteStringWith20PlaceHolder-4          1256984               930.8 ns/op           920 B/op         22 allocs/op
Benchmark_easytmpl_ExecuteStringWith20PlaceHolder-4            1870140               639.3 ns/op           880 B/op          3 allocs/op
Benchmark_FastTemplate_ExecuteFuncWith20PlaceHolder-4            1893668               618.1 ns/op           208 B/op         18 allocs/op
Benchmark_easytmpl_ExecuteFuncWith20PlaceHolder-4              1864115               615.9 ns/op           208 B/op         18 allocs/op
PASS
ok      github.com/tylitianrui/easytmpl/timing        9.133s

```
### benchmark results of  30 space holders 

```
➜  timing git:(master) ✗ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/tylitianrui/easytmpl/timing
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
Benchmark_GoTextTemplate_ExecuteWith30PlaceHolder-4               210078              5882 ns/op            1296 B/op         56 allocs/op
Benchmark_FastTemplate_ExecuteStringWith30PlaceHolder-4           723781              1556 ns/op            1623 B/op         32 allocs/op
Benchmark_easytmpl_ExecuteStringWith30PlaceHolder-4            1234134               962.6 ns/op          1328 B/op          3 allocs/op
Benchmark_FastTemplate_ExecuteFuncWith30PlaceHolder-4            1000000              1024 ns/op             312 B/op         27 allocs/op
Benchmark_easytmpl_ExecuteFuncWith30PlaceHolder-4              1192578              1122 ns/op             312 B/op         27 allocs/op
PASS
ok      github.com/tylitianrui/easytmpl/timing        9.353s
```


## 基本用法

### 非严格模式，不设置自动填充

```go
package main

import (
	"fmt"

	"github.com/tylitianrui/easytmpl"
)

func main() {
	tpl := "https://{{demain}}.com?name={{name}}&age={{age}}&birth={{birth}}"

	// Create a new template instance with default tag pair `{{` &` }}` and pre-allocated memory of 1024 bytes.
	t, err := easytmpl.NewTemplate(tpl,
		easytmpl.WithPreAllocateMemory(1024), // Pre-allocate memory for better performance ,1024 bytes
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

```

### 非严格模式,自动填充

```go
package main

import (
	"fmt"

	"github.com/tylitianrui/easytmpl"
)

func main() {
	tpl := "https://[[demain]].com?name=[[name]]&age=[[age]]&birth=[[birth]]"
	t, err := easytmpl.NewTemplate(tpl,
		easytmpl.WithTagPair("[[", "]]"),     // set custom tag pair `[[` & `]]`
		easytmpl.WithPreAllocateMemory(1024), // Pre-allocate memory for better performance ,1024 bytes
		easytmpl.WithAutoFill(""),            // Auto fill missing parameters with empty string
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

```

### 严格模式

```go
package main

import (
	"fmt"

	"github.com/tylitianrui/easytmpl"
)

func main() {
	tpl := "https://{{demain}}.com?name={{name}}&age={{age}}&birth={{birth}}"
	t, err := easytmpl.NewTemplate(tpl,
		easytmpl.WithTagPair("{{", "}}"),
		easytmpl.WithPreAllocateMemory(1024),
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

```


## 高级用法

```go
package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tylitianrui/easytmpl"
)

func main() {
	tpl := "https://{{demain}}.com?name={{name}}&age={{age}}&birth={{birth}}"
	t, err := easytmpl.NewTemplate(tpl)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = t.ExecuteFunc(&buf, func(w io.Writer, key string) (int, error) {
		switch key {
		case "demain":
			return w.Write([]byte("user.google"))
		case "name":
			return w.Write([]byte("tyltr"))
		case "age":
			return w.Write([]byte("18"))

		default:
			return w.Write([]byte("<null>")) //  Auto fill missing parameters with "<null>"
		}
	})
	fmt.Println("err:", err)
	fmt.Println("template:", buf.String())

	// Output:
	//err: <nil>
	//template: https://user.google.com?name=tyltr&age=18&birth=<null>

}

```



## 特殊案例说明
### 占位符为空，被当作普通文本


| 案例   | 模版    | 标签      |参数(*map类型*)     |效果    |
| :-----| :---- | :---- |:---- |:---- |
| 案例1 | `hello,i am {{name}}` | `{{`和 `}}`  |  `name`:`tyltr` |`hello,i am tyltr` |
| 案例2 | hello,i am **{{}}** {{name}} | `{{`和 `}}`  |  `name`:`tyltr` |`hello,i am {{}} tyltr` |

说明：案例2中包含空占位符 **{{}}** 会当作普通文本展示。


### 非贪婪模式匹配占位符

| 案例   | 模版    | 标签    |占位符    |参数(*map类型*)     |效果    |
| :-----| :---- | :---- |:----|:---- |:---- |
| 案例3 | `hello,i am {{{name}}` | `{{`和 `}}`  | **`name`** 而不是`{name`|  `name`:`tyltr` |`hello,i am {tyltr` |

说明：案例3 匹配采用非贪婪模式，会匹配尽可能少的字符。所以匹配出的占位符为 **`name`** 而不是`{name`

### 最左侧匹配原则

| 案例   | 模版    | 标签    |占位符    |参数(*map类型*)     |效果    |
| :-----| :---- | :---- |:----|:---- |:---- |
| 案例4 | `hello,i am {{{name{{}}tyltr}}` | `{{`和 `}}`  | **`name{{`** 而不是`}}tyltr`|  `name{{`:`ty` |`hello,i am {tytyltr` |

说明：
- `{{{name`**`{{}}`**`tyltr}}`  中间占位符为空，所以 **{{}}** 被当做普通文本处理
- 根据最左侧匹配原则  左侧 **`{{{name{{}}`** 拥有比右侧  **`{{}}tyltr}}`** 更高优先级
- 根据非贪婪原则，左侧部分会如此匹配  `{{{`**`name{{`**`}}` 

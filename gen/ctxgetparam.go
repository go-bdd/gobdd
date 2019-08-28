//+build ignore

package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type typeDef struct {
	Name  string
	Value string
	Bytes int
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func main() {
	funcMap := template.FuncMap{
		"Title":    strings.Title,
		"noescape": noescape,
	}

	floatTypes := []typeDef{
		{
			Name:  "float32",
			Value: "123.5",
			Bytes: 32,
		},
		{
			Name:  "float64",
			Value: "123.5",
			Bytes: 64,
		},
	}

	types := []typeDef{
		{
			Name:  "int",
			Value: "123",
			Bytes: 32,
		},
		{
			Name:  "int8",
			Value: "123",
			Bytes: 8,
		},
		{
			Name:  "int16",
			Value: "123",
			Bytes: 16,
		},
		{
			Name:  "int32",
			Value: "123",
			Bytes: 32,
		},
		{
			Name:  "int64",
			Value: "123",
			Bytes: 64,
		},
	}

	f, err := os.Create("context/getparam.go")
	die(err)

	var tmpl = template.Must(template.New("").Funcs(funcMap).Parse(getTmpl))

	err = tmpl.Execute(f, struct {
		IntTypes   []typeDef
		FloatTypes []typeDef
	}{IntTypes: types, FloatTypes: floatTypes})
	die(err)
	f.Close()

	f, err = os.Create("context/getparam_test.go")
	die(err)
	tmpl = template.Must(template.New("getparams").Funcs(funcMap).Parse(testTmpl))
	err = tmpl.Execute(f, struct {
		IntTypes   []typeDef
		FloatTypes []typeDef
	}{IntTypes: types, FloatTypes: floatTypes})
	die(err)
	f.Close()
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var testTmpl = `// Code generated .* DO NOT EDIT.
package context

import (
	"strconv"
	"testing"
	"fmt"
)
{{ range .FloatTypes }}
func TestContext_Get{{ .Name | Title }}Param(t *testing.T) {
	ctx := New()
	given := {{ .Name }}({{ .Value | noescape }})
	bs := float64ToByte(float64(given))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.Get{{ .Name | Title }}Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}
{{ end }}
func float64ToByte(f float64) []byte {
	return []byte(fmt.Sprintf("%f", f))
}
{{ range .IntTypes }}
func TestContext_Get{{ .Name | Title }}Param(t *testing.T) {
	ctx := New()
	given := {{ .Name }}(123)
	bs := []byte(strconv.Itoa(int(given)))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.Get{{ .Name | Title }}Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_Get{{ .Name | Title }}Param_PanicsOnNoSuchParam(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the Get{{ .Name | Title }}Param should panic")
        }
    }()
	_ = ctx.Get{{ .Name | Title }}Param(0)
}
{{ end }}
`

var getTmpl = `// Code generated .* DO NOT EDIT.
package context

import (
	"fmt"
	"strconv"
)

func (ctx Context) GetBoolParam(i int) bool {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	v, err := strconv.ParseBool(string(ctx.params[i]))
	if err != nil {
		panic(fmt.Sprintf("cannot convert to bool"))
	}

	return v
}

func (ctx Context) GetStringParam(i int) string {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	return string(ctx.params[i])
}

{{ range .FloatTypes }}
func (ctx Context) Get{{ .Name | Title }}Param(i int) {{ .Name }} {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseFloat(string(data), {{ .Bytes }})

	if err != nil {
		panic(err)
	}

	return {{ .Name }}(param)
}
{{ end }}
{{ range .IntTypes }}
func (ctx Context) Get{{ .Name | Title }}Param(i int) {{ .Name }} {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, {{ .Bytes }})

	if err != nil {
		panic(err)
	}

	return {{ .Name }}(param)
}
{{ end }}
`

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
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func main() {
	funcMap := template.FuncMap{
		"Title":    strings.Title,
		"noescape": noescape,
	}

	types := []typeDef{
		{
			Name:  "string",
			Value: `"example text"`,
		},
		{
			Name:  "int",
			Value: "123",
		},
		{
			Name:  "int8",
			Value: "123",
		},
		{
			Name:  "int16",
			Value: "123",
		},
		{
			Name:  "int32",
			Value: "123",
		},
		{
			Name:  "int64",
			Value: "123",
		},
		{
			Name:  "float32",
			Value: "123.5",
		},
		{
			Name:  "float64",
			Value: "123.5",
		},
		{
			Name:  "bool",
			Value: "false",
		},
	}

	f, err := os.Create("context/get.go")
	die(err)

	var tmpl = template.Must(template.New("").Funcs(funcMap).Parse(getTmpl))

	err = tmpl.Execute(f, struct {
		Types []typeDef
	}{Types: types})
	die(err)
	f.Close()

	f, err = os.Create("context/get_test.go")
	die(err)
	tmpl = template.Must(template.New("").Funcs(funcMap).Parse(testTmpl))
	err = tmpl.Execute(f, struct {
		Types []typeDef
	}{Types: types})
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

import "testing"
import "errors"

func TestContext_GetError(t *testing.T) {
	ctx := New()
	expected := errors.New("new err")
	ctx.Set("test", expected)
	received := ctx.GetError("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

{{ range .Types }}
func TestContext_Get{{ .Name | Title }}(t *testing.T) {
	ctx := New()
	expected := {{ .Name }}({{ .Value | noescape }})
	ctx.Set("test", expected)
	received := ctx.Get{{ .Name | Title }}("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_Get{{ .Name | Title }}_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := {{ .Name }}({{ .Value | noescape }})
	received := ctx.Get{{ .Name | Title }}("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_Get{{ .Name | Title }}_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the Get{{ .Name | Title }} should panic")
        }
    }()
	_ = ctx.Get{{ .Name | Title }}("test", {{ .Value | noescape }}, {{ .Value | noescape }})
}

func TestContext_Get{{ .Name | Title }}_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the Get{{ .Name | Title }} should panic")
        }
    }()
	_ = ctx.Get{{ .Name | Title }}("test")
}
{{ end }}	
`

var getTmpl = `// Code generated .* DO NOT EDIT.	
package context

import "fmt"

{{ range .Types }}
func (ctx Context) Get{{ .Name | Title }}(key interface{}, defaultValue ...{{ .Name }}) {{ .Name }} {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].({{ .Name }})
	if !ok {
		panic(fmt.Sprintf("the expected value is not {{ .Name }} (%T)", key))
	}
	return value
}
{{ end }}
`

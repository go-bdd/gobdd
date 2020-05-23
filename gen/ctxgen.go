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
	Zero  string
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
			Zero:  `""`,
		},
		{
			Name:  "int",
			Value: "123",
			Zero:  "0",
		},
		{
			Name:  "int8",
			Value: "123",
			Zero:  "0",
		},
		{
			Name:  "int16",
			Value: "123",
			Zero:  "0",
		},
		{
			Name:  "int32",
			Value: "123",
			Zero:  "0",
		},
		{
			Name:  "int64",
			Value: "123",
			Zero:  "0",
		},
		{
			Name:  "float32",
			Value: "123.5",
			Zero:  "0",
		},
		{
			Name:  "float64",
			Value: "123.5",
			Zero:  "0",
		},
		{
			Name:  "bool",
			Value: "false",
			Zero:  "false",
		},
	}

	f, err := os.Create("context_get.go")
	die(err)

	var tmpl = template.Must(template.New("").Funcs(funcMap).Parse(getTmpl))

	err = tmpl.Execute(f, struct {
		Types []typeDef
	}{Types: types})
	die(err)
	f.Close()

	f, err = os.Create("context_get_test.go")
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

const testTmpl = `// Code generated .* DO NOT EDIT.	
package gobdd

import "testing"
import "errors"

func TestContext_GetError(t *testing.T) {
	ctx := NewContext()
	expected := errors.New("new err")
	ctx.Set("test", expected)
	received, err := ctx.GetError("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

{{ range .Types }}
func TestContext_Get{{ .Name | Title }}(t *testing.T) {
	ctx := NewContext()
	expected := {{ .Name }}({{ .Value | noescape }})
	ctx.Set("test", expected)
	received, err := ctx.Get{{ .Name | Title }}("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_Get{{ .Name | Title }}_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := {{ .Name }}({{ .Value | noescape }})
	received, err := ctx.Get{{ .Name | Title }}("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_Get{{ .Name | Title }}_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.Get{{ .Name | Title }}("test", {{ .Value | noescape }}, {{ .Value | noescape }})
	if err == nil  {
		t.Error("the Get{{ .Name | Title }} should return an error")
	}
}

func TestContext_Get{{ .Name | Title }}_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.Get{{ .Name | Title }}("test")
	if err == nil  {
		t.Error("the Get{{ .Name | Title }} should return an error")
	}
}
{{ end }}	
`

const getTmpl = `// Code generated .* DO NOT EDIT.	
package gobdd

import "fmt"

{{ range .Types }}
func (ctx Context) Get{{ .Name | Title }}(key interface{}, defaultValue ...{{ .Name }}) ({{ .Name }}, error) {
	if len(defaultValue) > 1 {
        return {{.Zero|noescape}}, fmt.Errorf("allowed to pass only 1 default value but %d got", len(defaultValue))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0], nil
		}
		return {{.Zero|noescape}}, fmt.Errorf("the key %+v does not exist", key)
	}

	value, ok := ctx.values[key].({{ .Name }})
	if !ok {
		return {{.Zero|noescape}}, fmt.Errorf("the expected value is not {{ .Name }} (%T)", key)
	}
	return value, nil
}
{{ end }}
`

package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Create("context/getparams.go")
	die(err)
	defer f.Close()

	funcMap := template.FuncMap{
		"Title": strings.Title,
	}

	types := []string{
		"string",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
		"bool",
	}

	var tmpl = template.Must(template.New("").Funcs(funcMap).Parse(getTmpl))

	err = tmpl.Execute(f, struct {
		Types []string
	}{Types: types})
	die(err)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var getTmpl = `// Code generated .* DO NOT EDIT.
package context

import "fmt"
{{ range .Types }}
func (ctx Context) Get{{ . | Title }}(key interface{}, defaultValue ...{{ . }}) {{ . }} {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }
	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	value, ok := ctx.values[key].({{ . }})
	if !ok {
		panic(fmt.Sprintf("the expected value is not {{ . }} (%T)", key))
	}

	return value
}
{{ end }}
`

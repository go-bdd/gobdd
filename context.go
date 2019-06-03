package gobdd

import (
	"fmt"
	"strconv"
)

type Context struct {
	values map[string]interface{}
	params [][]byte
}

func newContext() Context  {
	return Context{
		values: map[string]interface{}{},
		params: [][]byte{},
	}
}

func (ctx Context) GetParam(i int) interface{} {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	return ctx.params[i]
}

func (ctx Context) GetStringParam(i int) string {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	return string(ctx.params[i])
}

func (ctx Context) GetFloat64Param(i int) float64 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseFloat(string(data), 32)

	if err != nil {
		panic(err)
	}

	return param
}

func (ctx Context) GetFloat32Param(i int) float32 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseFloat(string(data), 32)

	if err != nil {
		panic(err)
	}

	return float32(param)
}

func (ctx Context) GetIntParam(i int) int {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.Atoi(string(data))

	if err != nil {
		panic(err)
	}

	return param
}

func (ctx Context) Set(key string, value interface{}) {
	ctx.values[key] = value
}

func (ctx Context) Get(key string) interface{} {
	if _, ok := ctx.values[key]; !ok {
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	return ctx.values[key]
}

func (ctx Context) GetString(key string) string {
	if _, ok := ctx.values[key]; !ok {
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	value, ok := ctx.values[key].(string)
	if !ok {
		panic(fmt.Sprintf("the expected value is not string (%T)", key))
	}

	return value
}

func (ctx Context) GetFloat64(key string) float64 {
	if _, ok := ctx.values[key]; !ok {
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	value, ok := ctx.values[key].(float64)
	if !ok {
		panic(fmt.Sprintf("the expected value is not float64 (%T)", key))
	}

	return value
}

func (ctx Context) GetFloat32(key string) float32 {
	if _, ok := ctx.values[key]; !ok {
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	value, ok := ctx.values[key].(float32)
	if !ok {
		panic(fmt.Sprintf("the expected value is not float32 (%T)", key))
	}

	return value
}

func (ctx Context) GetInt(key string) int {
	if _, ok := ctx.values[key]; !ok {
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	value, ok := ctx.values[key].(int)
	if !ok {
		panic(fmt.Sprintf("the expected value is not int (%T)", key))
	}

	return value
}

func (ctx *Context) setParams(params [][]byte) {
	ctx.params = params
}

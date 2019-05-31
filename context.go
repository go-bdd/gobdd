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

package context

import (
	"fmt"
	"strconv"
)

// Holds two types of data:
//
// * data saved by previously executed steps
// * parameters which were received from the step definition
type Context struct {
	values map[interface{}]interface{}
	params [][]byte
}

func New() Context {
	return Context{
		values: map[interface{}]interface{}{},
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

func (ctx Context) GetInt8Param(i int) int8 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, 8)

	if err != nil {
		panic(err)
	}

	return int8(param)
}

func (ctx Context) GetInt16Param(i int) int16 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, 16)

	if err != nil {
		panic(err)
	}

	return int16(param)
}

func (ctx Context) GetInt32Param(i int) int32 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, 32)

	if err != nil {
		panic(err)
	}

	return int32(param)
}

func (ctx Context) GetInt64Param(i int) int64 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %d does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, 64)

	if err != nil {
		panic(err)
	}

	return param
}

func (ctx Context) Set(key interface{}, value interface{}) {
	ctx.values[key] = value
}

func (ctx Context) Get(key interface{}, defaultValue ...interface{}) interface{} {
	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %s does not exist", key))
	}

	return ctx.values[key]
}

func (ctx *Context) SetParams(params [][]byte) {
	ctx.params = params
}

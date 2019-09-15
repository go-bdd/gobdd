// Code generated .* DO NOT EDIT.
package context

import (
	"fmt"
	"strconv"
)

func (ctx Context) GetBoolParam(i int) bool {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
	}

	v, err := strconv.ParseBool(string(ctx.params[i]))
	if err != nil {
		panic(fmt.Sprintf("cannot convert to bool"))
	}

	return v
}

func (ctx Context) GetStringParam(i int) string {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
	}

	return string(ctx.params[i])
}


func (ctx Context) GetFloat32Param(i int) float32 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseFloat(string(data), 32)

	if err != nil {
		panic(err)
	}

	return float32(param)
}

func (ctx Context) GetFloat64Param(i int) float64 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseFloat(string(data), 64)

	if err != nil {
		panic(err)
	}

	return float64(param)
}


func (ctx Context) GetIntParam(i int) int {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, 32)

	if err != nil {
		panic(err)
	}

	return int(param)
}

func (ctx Context) GetInt8Param(i int) int8 {
	if i >= len(ctx.params) {
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
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
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
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
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
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
		panic(fmt.Sprintf("the param with index %+v does not exist", i))
	}

	data := ctx.params[i]
	param, err := strconv.ParseInt(string(data), 10, 64)

	if err != nil {
		panic(err)
	}

	return int64(param)
}


package context

import (
	"fmt"
)

// Holds two types of data:
//
// * data saved by previously executed steps
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

func (ctx Context) Set(key interface{}, value interface{}) {
	ctx.values[key] = value
}

func (ctx Context) Get(key interface{}, defaultValue ...interface{}) interface{} {
	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	return ctx.values[key]
}

package context

import (
	"fmt"
)

// Holds two types of data:
//
// * data saved by previously executed steps
type Context struct {
	values map[interface{}]interface{}
}

func New() Context {
	return Context{
		values: map[interface{}]interface{}{},
	}
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

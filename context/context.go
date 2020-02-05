package context

import (
	"fmt"
)

// Holds data from previously executed steps
type Context struct {
	values map[interface{}]interface{}
}

// Creates a new (empty) context struct
func New() Context {
	return Context{
		values: map[interface{}]interface{}{},
	}
}

// Sets the value under the key
func (ctx Context) Set(key interface{}, value interface{}) {
	ctx.values[key] = value
}

// Returns the data under the key.
// If couldn't find anything but the default value is provided, returns the default value.
// Otherwise, it panics.
func (ctx Context) Get(key interface{}, defaultValue ...interface{}) interface{} {
	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	return ctx.values[key]
}

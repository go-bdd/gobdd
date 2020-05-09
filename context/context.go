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
// Otherwise, it returns an error.
func (ctx Context) Get(key interface{}, defaultValue ...interface{}) (interface{}, error) {
	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0], nil
		}
		return nil, fmt.Errorf("the key %+v does not exist", key)
	}

	return ctx.values[key], nil
}

// It is a shortcut for getting the value already casted as error.
func (ctx Context) GetError(key interface{}, defaultValue ...interface{}) (interface{}, error) {
	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0], nil
		}
		return nil, fmt.Errorf("the key %+v does not exist", key)
	}

	if ctx.values[key] == nil {
		return nil, nil
	}

	value, ok := ctx.values[key].(error)
	if !ok {
		return nil, fmt.Errorf("the expected value is not error  (%T)", key)
	}
	return value, nil
}

package gobdd

import (
	"context"
	"errors"
	"reflect"
)

func validateStepFunc(f interface{}) error {
	value := reflect.ValueOf(f)
	if value.Kind() != reflect.Func {
		return errors.New("the parameter should be a function")
	}

	if value.Type().NumOut() != 2 {
		return errors.New("the function should return the context.Context and error")
	}
	val := value.Type().Out(0)
	contextInterface  := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !val.Implements(contextInterface) {
		return errors.New("the returned value should implement the context.Context interface")
	}

	val = value.Type().Out(1)
	errorInterface  := reflect.TypeOf((*error)(nil)).Elem()
	if !val.Implements(errorInterface) {
		return errors.New("the returned value should implement the Error interface")
	}

	if value.Type().NumIn() < 1 {
		return errors.New("the function should have Context as the first argument")
	}

	val = value.Type().In(0)
	n := val.ConvertibleTo(reflect.TypeOf((*context.Context)(nil)).Elem())
	if !n {
		return errors.New("the returned value should implement the Error interface")
	}
	return nil
}

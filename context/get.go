// Code generated .* DO NOT EDIT.	
package context

import "fmt"

func (ctx Context) GetString(key interface{}, defaultValue ...string) string {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(string)
	if !ok {
		panic(fmt.Sprintf("the expected value is not string (%T)", key))
	}
	return value
}

func (ctx Context) GetInt(key interface{}, defaultValue ...int) int {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(int)
	if !ok {
		panic(fmt.Sprintf("the expected value is not int (%T)", key))
	}
	return value
}

func (ctx Context) GetInt8(key interface{}, defaultValue ...int8) int8 {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(int8)
	if !ok {
		panic(fmt.Sprintf("the expected value is not int8 (%T)", key))
	}
	return value
}

func (ctx Context) GetInt16(key interface{}, defaultValue ...int16) int16 {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(int16)
	if !ok {
		panic(fmt.Sprintf("the expected value is not int16 (%T)", key))
	}
	return value
}

func (ctx Context) GetInt32(key interface{}, defaultValue ...int32) int32 {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(int32)
	if !ok {
		panic(fmt.Sprintf("the expected value is not int32 (%T)", key))
	}
	return value
}

func (ctx Context) GetInt64(key interface{}, defaultValue ...int64) int64 {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(int64)
	if !ok {
		panic(fmt.Sprintf("the expected value is not int64 (%T)", key))
	}
	return value
}

func (ctx Context) GetFloat32(key interface{}, defaultValue ...float32) float32 {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(float32)
	if !ok {
		panic(fmt.Sprintf("the expected value is not float32 (%T)", key))
	}
	return value
}

func (ctx Context) GetFloat64(key interface{}, defaultValue ...float64) float64 {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(float64)
	if !ok {
		panic(fmt.Sprintf("the expected value is not float64 (%T)", key))
	}
	return value
}

func (ctx Context) GetBool(key interface{}, defaultValue ...bool) bool {
	if len(defaultValue) > 1 {
        panic(fmt.Sprintf("allowed to pass only 1 default value but %d got", len(defaultValue)))
    }

	if _, ok := ctx.values[key]; !ok {
		if len(defaultValue) == 1 {
			return defaultValue[0]
		}
		panic(fmt.Sprintf("the key %+v does not exist", key))
	}

	value, ok := ctx.values[key].(bool)
	if !ok {
		panic(fmt.Sprintf("the expected value is not bool (%T)", key))
	}
	return value
}


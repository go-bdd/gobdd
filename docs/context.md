---
layout: default
title: Context
---

# Context

The context contains the data (state) from previous steps.

#### Passing data between steps

The context holds all the data from previously executed steps. They are accessible by `Context.GetX(key interface{})` functions:

* `Context.GetInt(key interface{}) int`
* `Context.GetFloat32(key interface{}) float32`
* `Context.GetFloat64(key interface{}) float64`
* `Context.GetString(key interface{}) string`
* and so on...

When you want to share some data between steps, use the `Context.Set(key, value interface{})` function

```go
// in the first step
ctx.Set(name{}, "John")

// in the second step
fmt.Printf("Hi %s\n", ctx.GetString(name{})) // prints "Hi John"
```

When the data is not provided, the whole test will fail.

## Good practices

It's a good practice to use custom structs as keys instead of strings or any built-in types to avoid collisions between steps using context.

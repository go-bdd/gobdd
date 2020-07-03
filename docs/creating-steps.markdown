---
layout: default
title: Creating steps
---

# Creating steps

Every step function should accept the `StepTest` as the first 2 parameters and returns the `Context`. Here's an example:

```go
type StepFunc func(gobdd.StepTest, ctx gobdd.Context, var1 int, var2 string)
```

What's important to stress - the context is a [custom struct](https://github.com/go-bdd/gobdd/tree/master/context), not the built-in interface.
To retrieve information from previously executed you should use functions `ctx.Get*(0)`. Replace the `*` with the type you need. Examples:

* `ctx.GetInt(type1{})`
* `ctx.GetString(type2{})`
* `ctx.Get(type3{})` - returns raw `[]byte`

If the value does not exist the test will fail.

There's possibility to pass a default value as well.

```go
ctx.GetFloat32(myFloatValue{}, 123)
```

If the `myFloatValue{}` value doesn't exists the `123` will be returned.

## Hooks

There's a possibility to define hooks which might be helpful building useful reporting, visualization, etc.

* `WithBeforeStep(f func(ctx Context))` configures functions that should be executed before every step
* `WithAfterStep(f func(ctx Context))` configures functions that should be executed after every step

```go
suite := NewSuite(
    t,
    WithBeforeStep(func(ctx Context) {
        ctx.Set(time.Time{}, time.Now())
    }),
    WithAfterStep(func(ctx Context) {
        start, _ := ctx.Get(time.Time{})
        log.Printf("step took %s", time.Since(start.(time.Time)))
    }),
)
```

## Good practices

Steps should be immutable and only communicate through [the context]({{ site.baseurl }}/context.html).

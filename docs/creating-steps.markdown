---
layout: default
title: Creating steps
---

# Creating steps

Every step function should accept the `StepTest` as the first 2 parameters and returns the `Context`. Here's an example:

```go
type StepFunc func(gobdd.StepTest, ctx context.Context, var1 int, var2 string) context.Context
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

## Good practices

Steps should be immutable and only communicate through [the context]({{ site.baseurl }}/context.html).

---
layout: default
title: Context
---

# Context

The context contains the data (state) from previous steps.

#### Passing data between steps

The context holds all the data from previously executed steps. They are accessible by `Context.GetX(key interface{})` functions:

* `Context.GetInt(key interface{}) (int, error)`
* `Context.GetFloat32(key interface{}) (float32, error)`
* `Context.GetFloat64(key interface{}) (float64, error)`
* `Context.GetString(key interface{}) (string, error)`
* and so on...

When you want to share some data between steps, use the `Context.Set(key, value interface{})` function

```go
// in the first step
ctx.Set(name{}, "John")

// in the second step
val, err := ctx.GetString(name{})
fmt.Printf("Hi %s\n", val) // prints "Hi John"
```

When the data is not provided, the whole test will fail.

#### Predefined keys

The context holds current test state `testing.T`. It is accessible by calling `Context.Get(TestingTKey{})`. This is useful if you need access to the test state from scenario or step hooks.

It is also possible to access references to current feature and scenario by calling `Context.Get(FeatureNameKey{})` and `Context.Get(ScenarioNameKey{})` respectively.

```go
value, err := ctx.Get(FeatureKey{})
feature, ok := value.(*msgs.GherkinDocument_Feature)
```

```go
value, err := ctx.Get(ScenarioKey{})
scenario, ok := value.(*msgs.GherkinDocument_Feature_Scenario)
```

## Good practices

It's a good practice to use custom structs as keys instead of strings or any built-in types to avoid collisions between steps using context.

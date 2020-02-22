---
layout: default
title: Suite's options
---

# Suite's options

The suite can be confiugred using one of these functions:

* `RunInParallel()` - enables running steps in parallel. It uses the stanard `T.Parallel` function.
* `WithFeaturesPath(path string)` - configures the path where GoBDD should look for features. The default value is `features/*.feature`.
* `WithTags(tags []string)` - configures which tags should be run. Every tag has to start with `@`.
* `WithBeforeScenario(f func())` - this function `f` will be called before every scenario.
* `WithAfterScenario(f func())` - this funcion `f` will be called after every scenario.
* `WithIgnoredTags(tags []string)` - configures tags which should be ignored and excluded from execution.

## Usage

Here are some examples of the usage of those functions:

```go
suite := NewSuite(t, WithFeaturesPath("features/func_types.feature"))
```

```go
suite := NewSuite(t, WithFeaturesPath("features/tags.feature"), WithTags([]string{"@tag"}))
```


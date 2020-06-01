---
layout: default
title: Parameter types
---

# Parameter types

GoBDD has support for [parameter types](https://cucumber.io/docs/cucumber/cucumber-expressions/). There are a few predefined parameter types:

 * `{int}` - integer (-1 or 56)
 * `{float}` - float (0.4 or 234.4)
 * `{word}` - single word (`hello` or `pizza`)
 * `{text}` - single-quoted or double-quoted strings (`'I like pizza'` or `"I like pizza"`)

You can add your own parameter types using `AddParameterTypes()` function. Here are a few examples

```go
    s := gobdd.NewSuite(t)
	s.AddParameterTypes(`{int}`, []string{`(\d)`})
	s.AddParameterTypes(`{float}`, []string{`([-+]?\d*\.?\d*)`})
	s.AddParameterTypes(`{word}`, []string{`([\d\w]+)`})
	s.AddParameterTypes(`{text}`, []string{`"([\d\w\-\s]+)"`, `'([\d\w\-\s]+)'`})
```

The first argument accepts the parameter types. As the second parameter provides list of regular expressions that should replace the parameter.

Parameter types should be added Before adding any step.
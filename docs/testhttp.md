---
layout: default
title: testhttp package
---

# testhttp package

This package is designed to help you testing your handlers and the domain hidden behind it.
To be able to use it, you have to initialise using `testhttp.Setup()` function.

```golang
s := NewSuite(t, WithFeaturesPath("features/http.feature"))
router := http.NewServeMux()
testhttp.Build(s, router)

s.Run()
```

The router passed to the build function should already have routing defined to work correctly.
You have a set of pre-defined steps you can use out-of the box:

 * `I make a (GET|POST|PUT|DELETE|OPTIONS) request to "([^"]*)` - for making a request
 * `the response code equals (\d+)` - for checking HTTP status codes
 * `the response contains a valid JSON` - self explanatory :)
 * `the response is "(.*)"` testing the response body's content

 Here is an example scenario using these steps:

 ```gherkin
 Feature: testhttp
   Scenario: testing JSON validation
    When I make a GET request to "/json"
    Then the response contains a valid JSON
    And the response is "{"valid": "json"}"
```

More pre-defined steps are coming soon!

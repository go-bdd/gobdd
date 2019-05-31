# GOBDD

This is a BDD testing framework. Uses gherkin for the test's syntax.

## Usage

Add the package to your project:

```
go get github.com/bkielbasa/gobdd
```

Add a new test main_test.go:

```go
func TestScenarios(t *testing.T) {
	suite := NewSuite(t)
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`I the result should equal (\d+)`, check)
	suite.Run()
}
```

Inside `features` folder create your scenarios. Here is an example:

```gherkin
Feature: math operations
  Scenario: add two digits
    When I add 1 and 2
    Then I the result should equal 3
```

## Creating steps

Every step function should implement the `StepFunc` function.

```go
type StepFunc func(ctx Context) error
```

### Context

The context contains two kinds of information:
* the data from previous steps
* step's parameters fetched from the step's name


#### Passing data between steps

The context holds all the data from previously executed steps. They are accessible by `Context.GetX(key string)` functions:

* `Context.GetInt(key string) int`
* `Context.GetFloat32(key string) float32`
* `Context.GetFloat64(key string) float64`
* `Context.GetString(key string) string`
* and so on...

When you want to share some data between steps, use the `Context.Set(key string, value interface{})` function

```go
// in first step
ctx.Set("name", "John")

// in the second step
fmt.Printf("Hi %s\n", ctx.GetString("name")) // prints "Hi John"
```

When the data is not provided, the whole test will fail.

#### Getting data from steps

The context holds all the data about the step's parameters. The naming convention is similar to passing data between steps but `Param` suffix should be added:

* `Context.GetIntParam(key int) int`
* `Context.GetFloat32Param(key int) float32`
* `Context.GetFloat64Param(key int) float64`
* `Context.GetStringParam(key int) string`
* and so on...

When the data is not provided, the whole test will fail.

So, for the step below there are two params available:

```go
suite.AddStep(`I add (\d+) and (\d+)`, add)
```

* `ctx.GetIntParam(0)` ->



















---
layout: default
title: Quick start
---

Add the package to your project:

```
go get github.com/go-bdd/gobdd
```

Add a new test `main_test.go`:

```go
func add(t gobdd.StepTest, ctx context.Context, var1, var2 int) {
	res := var1 + var2
	ctx.Set("sumRes", res)
}

func check(t gobdd.StepTest, ctx context.Context, sum int) {
	received, err := ctx.GetInt("sumRes")
    if err != nil {
        t.Error(err)

        return
    }

	if sum != received {
        t.Error(errors.New("the math does not work for you"))
	}
}

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

and run tests

```bash
go test ./...
```

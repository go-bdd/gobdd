---
layout: home
---

## Why did I make the library?

There is [godog](https://github.com/DATA-DOG/godog) library for BDD tests in Go. I found this library useful but it run as an external application which compiles our code. It has a several disadvantages:

* no debugging (breakpoints) in the test. Sometimes it's useful to go through the whole execution step by step
* metrics don't count the test run this way
* some style checkers recognise tests as dead code
* it's impossible to use built-in features like [build constraints](https://golang.org/pkg/go/build/#hdr-Build_Constraints).
* no context in steps - so the state have to be stored somewhere else - in my opinion, it makes the maintenance harder

## More details about the implementation

The features use [gherkin](https://cucumber.io/docs/gherkin/reference/) syntax. It means that every document (which contains features etc) have to be compatalble with it.

Firstly, the library reads all available documents. By default, `features/*.feature` files. Then, loads all the step definitions. Next, tries to execute every scenario and steps into the scenario one by one. At the end, it produces the report of the execution.

You can have multiple gherkin documents executed within one test.

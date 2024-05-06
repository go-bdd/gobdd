
Feature: ignored tags
  @ignore
  Rule: this rule should be ignored
    Scenario: the scenario should be ignored
      Then fail the test
  Rule: this rule should run
    Scenario: the scenario should pass
      Then the test should pass

Feature: ignored tags
  @tag
  Scenario: the scenario should be pass
    Then the test should pass
  @tag
  Scenario Outline: the scenario should be pass
    Then the test should pass
    Examples:
  Scenario: the test should never be executed
    Then fail the test

  Rule: the rule should never be executed
    Scenario: the test in ignored rule should never be executed
      Then fail the test

  Rule: this rule should run
    @tag
    Scenario: the test in executed rule should pass
      Then the test should pass
    Scenario: the test in executed rule should never be executed
      Then fail the test
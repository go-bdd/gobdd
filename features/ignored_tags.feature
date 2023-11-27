Feature: ignored tags
  @ignore
  Scenario: the scenario should be ignored
    Then fail the test
  Scenario: the scenario should pass
    Then the test should pass

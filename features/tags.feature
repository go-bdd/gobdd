Feature: ignored tags
  @tag
  Scenario: the scenario should be pass
    Then the test should pass

  Scenario: the test should never be executed
    Then fail the test
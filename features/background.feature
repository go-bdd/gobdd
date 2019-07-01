Feature: using background steps
  Background: adding
    When I add 1 and 2

  Scenario: the background step should be executed
    Then the result should equal 3

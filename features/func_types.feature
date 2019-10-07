Feature: math operations in floats
  Scenario: add two floats
    When I add 1.0 and 2.3
    Then the result should equal 3.3
Feature: parameter types
  Scenario: add two digits
    When I add 1 and 2
    Then the result should equal 3
  Scenario: simple word
    When I use word pizza
  Scenario: simple text with double quotes
    When I use text "I like pizza"
  Scenario: simple text with single quotes
    When I use text 'I like pizza'
  Scenario: add two floats
    When I add floats 1 and 2
    Then the result should equal float 3
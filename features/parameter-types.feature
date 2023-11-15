Feature: parameter types
  Scenario: add two digits
    When I add 1 and 2
    Then the result should equal 3
  Scenario: simple word
    When I use word pizza
  Scenario: simple text with double quotes
    When I use text "I like pizza!"
    Then the result should equal text 'I like pizza!'
  Scenario: simple text with single quotes
    When I use text 'I like pizza!'
    Then the result should equal text 'I like pizza!'
  Scenario: add two floats
    When I add floats -1.2 and 2.4
    Then the result should equal float 1.2
  Scenario: concat a word and a text with single quotes
    When I concat word Hello and text ' World!'
    Then the result should equal text 'Hello World!'
  Scenario: concat a word and a text with double quotes
    When I concat word Hello and text " World!"
    Then the result should equal text "Hello World!"
  Scenario: format text
    When I format text "counter %d" with int -12
    Then the result should equal text "counter -12"

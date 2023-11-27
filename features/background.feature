Feature: using background steps
  Background: adding
    When I add 1 and 2

  Scenario: the background step should be executed
    Then the result should equal 3

  Rule: adding and concat
    Background: concat
      When I concat word Hello and text " World!"

    Scenario: the background steps should be executed
      Then the result should equal 3
      Then the result should equal text "Hello World!"

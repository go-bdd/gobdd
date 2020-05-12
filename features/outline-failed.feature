Feature: Scenario Outline
  Scenario Outline: testing outline scenarios
    When I add <digit1> and <digit2>
    Then the result should equal <result>
    Examples:
     | digit1 | digit2 | result |
     | 1 | 2 | 3                |
     | 5 | 5 | 11               |
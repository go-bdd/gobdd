Feature: Argument feature
  Scenario: compare text with argument
    When I concat text "Hello " and argument:
      """
      World!
      """
    Then the result should equal argument:
      """
      Hello World!
      """
  Scenario: compare text with multiline argument
    When I concat text "Hello " and argument:
      """
      New
      World!
      """
    Then the result should equal argument:
      """
      Hello New
      World!
      """

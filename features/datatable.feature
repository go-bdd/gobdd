Feature: dataTable feature
  Scenario: compare text with dataTable
    When I concat all the columns and row together using " - " to separate the columns
      |r1c1|r1c2|r1c3|
      |r2c1|r2c2|r2c3|
      |r3c1|r3c2|r3c3|
    Then the result should equal argument:
      """
      r1c1 - r1c2 - r1c3
      r2c1 - r2c2 - r2c3
      r3c1 - r3c2 - r3c3
      """
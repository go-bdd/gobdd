Feature: HTTP requests
  Scenario: test GET request
    When I make a GET request to "/health"
    Then the response code equals 200
  Scenario: not existing URI
    When I make a GET request to "/not-exists"
    Then the response code equals 404
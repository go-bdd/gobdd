Feature: HTTP requests
  Scenario: test GET request
    When I make a GET request to "/health"
    Then the response code equals 200
  Scenario: not existing URI
    When I make a GET request to "/not-exists"
    Then the response code equals 404
  Scenario: testing JSON validation
    When I make a GET request to "/json"
    Then the response contains a valid JSON
    And the response is "{"valid": "json"}"

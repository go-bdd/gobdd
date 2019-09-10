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

  Scenario: the simplest request
    When I have a GET request "http://google.com"
    Then the request has method set to GET
    And the url is set to "http://google.com"
    And the request body is nil
  Scenario: passing headers
    Given I have a GET request "http://google.com"
    When I set the header "XYZ" to "ZYX"
    Then the request has header "XYZ" set to "ZYX"
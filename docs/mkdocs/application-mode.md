# Application Mode

The application mode can be selected by manipulating the `X-App-Mode` header.

Docker compose has predefined endpoints and each endpoint hase another settings.

When you call `Geralt` by:

 - http://localhost:9000
 ```
 Default configuration.
 ```

 - http://geralt.localhost:9000
 ```
 X-App-Mode is set to "default"
 ```

- http://mock-geralt.localhost:9000
 ```
 X-App-Mode is set to "mock"
 ```

When `Geralt` is called in `"mock"` mode then can use the mock server without any calls to Cook API.
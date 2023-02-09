# Observability

You can connect to any OTEL agent you want. You can choose the one from [New Relic](https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/) or your local one.

You can setup the host where traces will be sent using env var `OTEL_EXPORTER_OTLP_ENDPOINT`. In production, it should probably point to `https://otlp.eu01.nr-data.net:443`.
If the agent requires authentication, you can provide the API KEY or barer token using the `OTEL_EXPORTER_OTLP_HEADERS` env var.

```
OTEL_EXPORTER_OTLP_HEADERS="api-key=YOUR_KEY_HERE"
```

For local testing, use the local OTEL agent with jaeger.

```sh
docker-compose -f compose-opentelemetry.yml up
```

And set env var to:

```
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318/
```
# scootgo

Generic Go service boilerplate with three runnable entry points:

- HTTP server via `serve`
- sample CLI command via `exec`
- sample RabbitMQ consumer via `consume`

## Project structure

- `main.go`: application bootstrap
- `pkg/server`: Cobra command wiring for `serve`, `exec`, and `consume`
- `pkg/app/handler`: sample HTTP endpoint
- `pkg/app/command`: sample CLI command
- `pkg/app/consumer`: sample RabbitMQ consumer
- `pkg/config`: environment-driven configuration
- `pkg/logger`: zerolog setup
- `pkg/rmq`: RabbitMQ connector, publisher, and consumer helpers
- `pkg/healthcheck`: health and readiness endpoints

## Run

Build the binary:

```bash
make build
```

Start the HTTP server:

```bash
make run-http
```

Default endpoint:

```bash
curl http://localhost:8080/api/v1/ping
```

Expected response:

```json
{"service":"scootgo","status":"ok"}
```

Run the sample command:

```bash
build/scootgo exec hello
```

Or pass a custom message:

```bash
build/scootgo exec hello "custom message"
```

Run the sample consumer:

```bash
make run-rmq-c QUEUE_NAME=sample ROUTING_KEY=sample.created
```

The sample consumer expects JSON payloads like:

```json
{"id":"123","message":"hello"}
```

## Configuration

Useful environment variables:

- `APP_ENV`: defaults to `prod`
- `APP_NAME`: defaults to `scootgo`
- `APP_PORT`: defaults to `8080`
- `SERVER_MODE`: defaults to `release`
- `LOG_LEVEL`: defaults to `1`
- `RABBIT_HOST`, `RABBIT_PORT`, `RABBIT_USER`, `RABBIT_PASSWORD`, `RABBIT_VIRT_HOST`: used by `consume`

`serve` and `exec` do not require DB or RabbitMQ to start.

## Development

Run tests:

```bash
make test
```

This runs the maintained boilerplate/unit packages.

Format code:

```bash
make fmt
```

Run linter:

```bash
make lint
```

## Extend the boilerplate

- Add new HTTP routes under `pkg/app/handler`
- Add new commands under `pkg/app/command`
- Add new consumers under `pkg/app/consumer`
- Register them in `pkg/server`

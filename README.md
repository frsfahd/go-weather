# go-weather

[Restful weather API](https://roadmap.sh/projects/weather-api-wrapper-service) based on [virtual crossing](https://www.visualcrossing.com/) service. This project adds cache layer to the based API.

Additionally, there is also performance testing with [k6](https://grafana.com/docs/k6/latest/) (see `load_test.js` and `smoke_test.js`).
For rate limiting, it uses token bucket from `golang.org/x/time/rate` (see `middlewares.go`)

## Getting Started

> Live demo : https://go-weather.fly.biz.id/docs

1. clone this repo
2. spin up redis instance `make docker-run`
3. create .env file (see `env-sample`)
4. run the API `make run`

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

Live reload the application:

```bash
make watch
```

Run the smoke test:

```bash
make smoke_test
```

Run the load test:

```bash
make load_test
```

Clean up binary from the last build:

```bash
make clean
```

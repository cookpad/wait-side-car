wait-side-car
=============

Wrapper command line tool to wait side-car Envoy containers on application startup.

- Wait Envoy side-car with retry every 300 millseconds.
  - Retry on connection errors in case Envoy hasn't listen the socket yet.
  - Retry on unhealth responses in case Envoy hasn't fetch upstream hosts/routes from xDS API yet.
  - If Envoy responds a healthy response, stop retrying.
- Timeout all waiting procedures after `--timeout` millseconds. Even if timed out, go next step (timeout is not failure here).
- Call execve(2) to start user program with current environment variables and rest of command line arguments.

## Usage
Common flags:

- `--timeout`: Overall timeout millseconds.
- `--envoy-host-header`: A host header to access to specific upstream serivce via Envoy.

For HTTP web APIs:

- `--envoy-healthcheck-url`: A URL to check accessibility to upstream service via Envoy (e.g. `http://envoy/hc` or `http://127.0.0.1:9211/hc`)

```
wait-side-car --timeout=10000 --envoy-host-header=user-service --envoy-healthcheck-url=http://envoy/hc /usr/local/bin/my-app -b 0.0.0.0 -p 8080 another-arg
```

For gRPC servers:

- `--envoy-grpc-insecure-healthcheck-addr`: A pair of IP address (or DNS name) and port to check accessibility to upstream gRPC service via Envoy (e.g `envoy:8080` or `127.0.0.1:9211`)

```
wait-side-car --timeout=10000 --envoy-host-header=user-service --envoy-grpc-insecure-healthcheck-addr=envoy:8080 /usr/local/bin/my-app run-job --job-name=record-count-of-users
```

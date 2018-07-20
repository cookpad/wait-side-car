wait-side-car
=============

Wrapper command line tool to wait side-car containers like Envoy on startup.

- Wait envoy side-car with retry every 300 millseconds.
  - Retry on connection errors in case Envoy hasn't listen the socket yet.
  - Retry on non-200 status codes in case Envoy hasn't fetch upstream hosts/routes from xDS API yet.
  - If Envoy responds 200, stop retrying.
- Timeout all waiting procedures after `--timeout` millseconds.
- Call execve(2) to start user program with current environment variables and rest of command line arguments.

## Usage
Typical use case:

- `--timeout`: Overall timeout millseconds.
- `--envoy-healthcheck-url`: A URL to check accessibility to upstream service via Envoy.
- `--envoy-host-header`: A host header to access to specific upstream serivce via Envoy.

```
wait-side-car --timeout=10000 --envoy-healthcheck-url=http://envoy/hello --envoy-host-header=xxx /usr/local/bin/my-app -b 0.0.0.0 -p 8080 another-arg
```

Or

```
wait-side-car -h
```

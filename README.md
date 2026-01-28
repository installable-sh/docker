# RUN & INSTALL

Minimal Docker image that fetches and executes shell scripts from URLs. Uses [mvdan/sh](https://github.com/mvdan/sh), a POSIX shell interpreter written in pure Go.

Published on Docker Hub as `installable/sh`.

## Usage

```dockerfile
FROM installable/sh AS installable

FROM ubuntu:latest
COPY --from=installable / /
```

The `COPY --from=installable / /` pattern adds the `RUN` and `INSTALL` binaries to any base image, allowing scripts to use the base image's utilities.

| Command                       | Description                                | Status |
| ----------------------------- | ------------------------------------------ | ------ |
| [`RUN`](#run-command)         | Fetch and execute scripts at runtime       | Ready  |
| [`INSTALL`](#install-command) | Installation and setup tasks during builds | WIP    |

See the [examples](./examples) directory for more examples.

---

## RUN Command

Fetch and execute scripts at runtime.

```dockerfile
CMD ["RUN", "https://example.com/script.sh", "arg1", "arg2"]
```

### Environment Variables as Headers

Use `+env` to send environment variables as HTTP headers when fetching scripts:

```dockerfile
CMD ["RUN", "+env", "https://example.com/script.sh"]
```

Each environment variable is sent as an `X-Env-*` header:

- `API_KEY=secret` → `X-Env-API_KEY: secret`
- `FOO=bar` → `X-Env-FOO: bar`

This allows dynamic script generation based on the container's environment.

### Custom User-Agent

The default User-Agent is `run/1.0 (installable)`. Set the `USER_AGENT` environment variable to override it:

```dockerfile
ENV USER_AGENT="MyApp/1.0"
CMD ["RUN", "https://example.com/script.sh"]
```

### Raw Output

Use `+raw` to print the fetched script without executing it:

```bash
RUN +raw https://example.com/script.sh
```

This is useful for debugging or piping the script to another tool.

### Bypass CDN Cache

Use `+nocache` to request fresh content from the origin server:

```bash
RUN +nocache https://example.com/script.sh
```

This sets `Cache-Control: no-cache, no-store, must-revalidate` and `Pragma: no-cache` headers.

### Features

- **Minimal footprint**: Small image with embedded CA certificates
- **No shell required**: The POSIX shell interpreter is embedded in the Go binary
- **Argument passing**: Pass arguments to scripts via `$1`, `$2`, etc.
- **HTTPS support**: CA certificates embedded in the binary
- **Environment forwarding**: Send env vars as HTTP headers with `+env`
- **Custom User-Agent**: Set `USER_AGENT` env var to customize the request header
- **Raw output**: Use `+raw` to print the script without executing
- **Cache bypass**: Use `+nocache` to skip CDN caches
- **Compatible**: Overlay onto any base image with `COPY --from=installable / /`

---

## INSTALL Command

🚧 **Work in Progress** 🚧

The `INSTALL` command is under development and will be available in a future release. It is intended for installation and setup tasks during Docker image builds.

---

## Development

### Quick Start

```bash
make        # Setup, format, lint, test, and build
```

### Makefile Targets

| Target         | Description                            |
| -------------- | -------------------------------------- |
| `make`         | Run all: setup, fmt, lint, test, build |
| `make setup`   | Install git hooks and golangci-lint    |
| `make fmt`     | Format code with gofmt                 |
| `make lint`    | Run golangci-lint                      |
| `make check`   | Check formatting (CI)                  |
| `make test`    | Run tests                              |
| `make RUN`     | Build RUN binary                       |
| `make INSTALL` | Build INSTALL binary                   |
| `make clean`   | Remove binaries and test cache         |

### Docker Build

```bash
docker build -t installable/sh .
```

**Note**: The Dockerfile copies CA certificates from Alpine's `ca-certificates` package during the build. The certificates in `internal/certs/` are placeholders for local development only.

## Limitations

- Scripts must be POSIX-compliant (no bash-specific features like arrays, `[[`, `set -E`, etc.)
- External commands must exist in the base image

## License

Copyright 2026 Scaffoldly LLC

Apache 2.0

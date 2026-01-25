# RUN

A minimal Docker image that fetches and executes shell scripts from URLs. Uses [mvdan/sh](https://github.com/mvdan/sh), a POSIX shell interpreter written in pure Go.

## Usage

```dockerfile
FROM scaffoldly/run AS run

FROM ubuntu:latest
COPY --from=run / /
CMD ["RUN", "https://example.com/script.sh", "arg1", "arg2"]
```

The `COPY --from=run / /` pattern adds the `RUN` binary to any base image, allowing scripts to use the base image's utilities.

Also available at `ghcr.io/scaffoldly/run`.

See the [examples](./examples) directory for more examples.

## Environment Variables as Headers

Use `+env` to send environment variables as HTTP headers when fetching scripts:

```dockerfile
CMD ["RUN", "+env", "https://example.com/script.sh"]
```

Each environment variable is sent as an `X-Env-*` header:
- `API_KEY=secret` → `X-Env-API_KEY: secret`
- `FOO=bar` → `X-Env-FOO: bar`

This allows dynamic script generation based on the container's environment.

## Features

- **Minimal footprint**: ~4MB image (compressed Go binary with embedded CA certificates)
- **No shell required**: The POSIX shell interpreter is embedded in the Go binary
- **Argument passing**: Pass arguments to scripts via `$1`, `$2`, etc.
- **HTTPS support**: CA certificates embedded in the binary
- **Environment forwarding**: Send env vars as HTTP headers with `+env`
- **Compatible**: Overlay onto any base image with `COPY --from=run / /`

## Building

```bash
docker build -t run .
```

## Limitations

- Scripts must be POSIX-compliant (no bash-specific features like arrays, `[[`, `set -E`, etc.)
- External commands must exist in the base image

## License

MIT

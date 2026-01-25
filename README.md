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

## Features

- **Minimal footprint**: ~9MB image (Go binary + CA certificates)
- **No shell required**: The POSIX shell interpreter is embedded in the Go binary
- **Argument passing**: Pass arguments to scripts via `$1`, `$2`, etc.
- **HTTPS support**: Includes CA certificates for fetching scripts over HTTPS
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

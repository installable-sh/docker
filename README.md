# RUN & INSTALL

Minimal Docker image that fetches and executes shell scripts from URLs. Uses [mvdan/sh](https://github.com/mvdan/sh), a POSIX shell interpreter written in pure Go.

Includes two commands:
- **RUN**: Fetch and execute scripts (for runtime execution)
- **INSTALL**: Coming soon (for installation/setup tasks)

## Usage

```dockerfile
FROM installable/sh AS run

FROM ubuntu:latest
COPY --from=run / /
CMD ["RUN", "https://example.com/script.sh", "arg1", "arg2"]
```

The `COPY --from=run / /` pattern adds the `RUN` and `INSTALL` binaries to any base image, allowing scripts to use the base image's utilities.

Published on Docker Hub as `installable/sh`.

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

## Custom User-Agent

The default User-Agent is `run/1.0 (installable)`. Set the `USER_AGENT` environment variable to override it:

```dockerfile
ENV USER_AGENT="MyApp/1.0"
CMD ["RUN", "https://example.com/script.sh"]
```

## Raw Output

Use `+raw` to print the fetched script without executing it:

```bash
RUN +raw https://example.com/script.sh
```

This is useful for debugging or piping the script to another tool.

## Bypass CDN Cache

Use `+nocache` to request fresh content from the origin server:

```bash
RUN +nocache https://example.com/script.sh
```

This sets `Cache-Control: no-cache, no-store, must-revalidate` and `Pragma: no-cache` headers.

## Features

- **Minimal footprint**: ~4MB image (compressed Go binary with embedded CA certificates)
- **No shell required**: The POSIX shell interpreter is embedded in the Go binary
- **Argument passing**: Pass arguments to scripts via `$1`, `$2`, etc.
- **HTTPS support**: CA certificates embedded in the binary
- **Environment forwarding**: Send env vars as HTTP headers with `+env`
- **Custom User-Agent**: Set `USER_AGENT` env var to customize the request header
- **Raw output**: Use `+raw` to print the script without executing
- **Cache bypass**: Use `+nocache` to skip CDN caches
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

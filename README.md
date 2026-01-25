# RUN

A minimal Docker image that fetches and executes shell scripts from URLs without requiring a shell interpreter to be installed. Uses [mvdan/sh](https://github.com/mvdan/sh), a POSIX shell interpreter written in pure Go.

## Usage

```dockerfile
FROM ghcr.io/scaffoldly/run
CMD ["https://example.com/script.sh", "arg1", "arg2"]
```

Or run directly:

```bash
docker run ghcr.io/scaffoldly/run https://example.com/script.sh arg1 arg2
```

## Features

- **Minimal footprint**: ~15MB image (Go binary + busybox + CA certificates)
- **No shell required**: The POSIX shell interpreter is embedded in the Go binary
- **Argument passing**: Pass arguments to scripts via `$1`, `$2`, etc.
- **HTTPS support**: Includes CA certificates for fetching scripts over HTTPS
- **Common utilities**: Includes busybox for `sleep`, `echo`, `cat`, `grep`, etc.

## Building

```bash
docker build -t run .
```

## Limitations

- Scripts must be POSIX-compliant (no bash-specific features like arrays, `[[`, `set -E`, etc.)
- External commands must exist in the image (busybox provides common ones)

## License

MIT

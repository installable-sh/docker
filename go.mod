module github.com/installable-sh/docker/v1

go 1.25.4

require (
	github.com/hashicorp/go-retryablehttp v0.7.8
	mvdan.cc/sh/v3 v3.12.0
)

require (
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
)

replace (
	github.com/installable-sh/docker/v1/internal/certs => ./internal/certs
	github.com/installable-sh/docker/v1/internal/fetch => ./internal/fetch
	github.com/installable-sh/docker/v1/internal/install => ./internal/install
	github.com/installable-sh/docker/v1/internal/run => ./internal/run
	github.com/installable-sh/docker/v1/internal/shell => ./internal/shell
	github.com/installable-sh/docker/v1/internal/version => ./internal/version
)

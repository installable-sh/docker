FROM golang:alpine AS builder
RUN apk add --no-cache git ca-certificates make && \
    apk add --no-cache upx || true
WORKDIR /build
COPY go.mod go.sum Makefile ./
RUN --mount=type=cache,target=/go go mod download
COPY . ./
# Copy Alpine's CA certificates to the location our code expects
RUN cp /etc/ssl/certs/ca-certificates.crt internal/certs/ca-certificates.crt
RUN --mount=type=cache,target=/go make all compress

FROM scratch AS combined
COPY --from=builder /build/RUN /usr/local/bin/RUN
COPY --from=builder /build/INSTALL /usr/local/bin/INSTALL
RUN ["/usr/local/bin/RUN", "--help"]
RUN ["/usr/local/bin/INSTALL", "--help"]

FROM scratch
COPY --from=combined / /
ENTRYPOINT ["RUN"]
CMD ["--help"]

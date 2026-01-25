FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o run .

FROM busybox AS combined
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/run /bin/RUN

FROM scratch
COPY --from=combined / /
ENTRYPOINT ["RUN"]
CMD ["--help"]

FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY hack/ca-certificates.crt ./hack/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o run . && \
    (apk add --no-cache upx && upx --best --lzma run || true)

FROM scratch
COPY --from=builder /build/run /usr/local/bin/RUN
ENTRYPOINT ["RUN"]
CMD ["--help"]

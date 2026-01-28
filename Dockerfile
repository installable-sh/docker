FROM golang:alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
# Copy Alpine's CA certificates to the location our code expects
RUN cp /etc/ssl/certs/ca-certificates.crt internal/certs/ca-certificates.crt
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o RUN ./cmd/run && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o INSTALL ./cmd/install && \
    (apk add --no-cache upx && upx --best --lzma RUN INSTALL || true)

FROM scratch AS combined
COPY --from=builder /build/RUN /usr/local/bin/RUN
COPY --from=builder /build/INSTALL /usr/local/bin/INSTALL
RUN ["/usr/local/bin/RUN", "--help"]
RUN ["/usr/local/bin/INSTALL", "--help"]

FROM scratch
COPY --from=combined / /
ENTRYPOINT ["RUN"]
CMD ["--help"]

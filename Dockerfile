FROM golang:alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o RUN ./cmd/run && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o INSTALL ./cmd/install && \
    (apk add --no-cache upx && upx --best --lzma RUN INSTALL || true)

FROM scratch
COPY --from=builder /build/RUN /usr/local/bin/RUN
COPY --from=builder /build/INSTALL /usr/local/bin/INSTALL
ENTRYPOINT ["RUN"]
CMD ["--help"]

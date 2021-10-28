FROM golang:latest AS builder

WORKDIR /build
COPY . /build
RUN go mod tidy
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/geoip2-server /
ENTRYPOINT ["/geoip2-server"]
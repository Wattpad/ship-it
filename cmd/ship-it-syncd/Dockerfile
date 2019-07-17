FROM golang:1.12 AS golang-build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o ship-it-syncd cmd/ship-it-syncd/main.go

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=golang-build /build/ship-it-syncd /
CMD ["/ship-it-syncd"]

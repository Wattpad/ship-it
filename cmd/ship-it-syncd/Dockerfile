FROM golang:1.13 AS golang-build
WORKDIR /build
COPY operator/go.mod operator/go.sum ./operator/
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ship-it-syncd/main.go cmd/ship-it-syncd/main.go
COPY internal ./internal/
COPY operator/api ./operator/api
RUN CGO_ENABLED=0 go build -o ship-it-syncd cmd/ship-it-syncd/main.go

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=golang-build /build/ship-it-syncd /
CMD ["/ship-it-syncd"]

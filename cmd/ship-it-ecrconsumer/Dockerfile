FROM golang:1.12 AS golang-build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o ship-it-ecrconsumer cmd/ship-it-ecrconsumer/main.go

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=golang-build /build/ship-it-ecrconsumer /
CMD ["/ship-it-ecrconsumer"]

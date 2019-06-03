FROM golang:1.12 AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o ship-it

FROM alpine:3.8
COPY --from=build /build/ship-it /
CMD ["/ship-it"]

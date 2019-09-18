FROM golang:1.13 AS golang-build
WORKDIR /build
COPY operator/go.mod operator/go.sum ./operator/
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ship-it-api/main.go cmd/ship-it-api/main.go
COPY internal/ ./internal/
COPY operator/api ./operator/api
RUN CGO_ENABLED=0 go build -o ship-it-api cmd/ship-it-api/main.go

FROM node:12.3 AS node-build
COPY web .
RUN npm install react-scripts && npm run build

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=golang-build /build/ship-it-api /
COPY --from=node-build /build /dashboard
CMD ["/ship-it-api"]

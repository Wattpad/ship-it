FROM golang:1.12 AS golang-build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o ship-it

FROM node:12.3 AS node-build
COPY web .
RUN npm install react-scripts && npm run build

FROM alpine:3.8
COPY --from=golang-build /build/ship-it /
COPY --from=node-build /build /dashboard
CMD ["/ship-it"]

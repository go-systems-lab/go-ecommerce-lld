FROM golang:1.24.3-alpine3.22 AS builder
RUN apk update && apk --no-cache add build-base ca-certificates
WORKDIR /go/src/github.com/go-systems-lab/go-ecommerce-lld
COPY go.mod go.sum ./
RUN go mod download
COPY account account
COPY product product
COPY order order
COPY tools.go .
RUN GO111MODULE=on go build -o /go/bin/app ./order/cmd/order

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=builder /go/bin/app .
EXPOSE 8080
CMD ["./app"]
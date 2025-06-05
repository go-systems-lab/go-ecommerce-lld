FROM golang:1.24.3-alpine3.22 AS builder
RUN apk update && apk --no-cache add build-base ca-certificates
WORKDIR /go/src/github.com/go-systems-lab/go-ecommerce-lld
COPY go.mod go.sum ./
RUN go mod download
COPY account account
COPY tools.go .
RUN GO111MODULE=on go build -o /go/bin/app ./account/cmd/account

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=builder /go/bin/app .
EXPOSE 50051
CMD ["./app"]
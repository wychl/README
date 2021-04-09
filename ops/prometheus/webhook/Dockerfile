FROM golang:1.12.6-alpine3.10 AS builder
WORKDIR /go/src/webhook
COPY . .
RUN go build -o server .

FROM alpine:3.10
ENV TZ Asia/Shanghai
WORKDIR /root/
COPY --from=builder /go/src/webhook/server /usr/local/bin
CMD ["server"]
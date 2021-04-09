FROM golang:1.10-alpine as builder
WORKDIR /go/src/app
COPY main.go .
RUN  go build -o server .

FROM envoyproxy/envoy-alpine-dev:latest
WORKDIR /usr/local/bin
COPY --from=builder /go/src/app/server .
COPY start_service.sh .
RUN chmod u+x start_service.sh
RUN chmod u+x server
ENTRYPOINT /usr/local/bin/start_service.sh

FROM golang:1.7.3 as builder
WORKDIR /go/src/app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /go/src/app/server .
EXPOSE 8888
CMD ["./server"]
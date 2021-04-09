FROM golang:alpine as builder
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"
WORKDIR /app 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM scratch
WORKDIR /app
COPY --from=builder /app/server /usr/bin/
ENTRYPOINT ["server"]
FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add -U --no-cache tzdata ca-certificates
ENV TZ Asia/Shanghai
WORKDIR /app
ADD server .
EXPOSE 8080
CMD [ "./server" ]
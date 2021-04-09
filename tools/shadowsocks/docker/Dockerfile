FROM ubuntu:latest
    
WORKDIR /app

COPY ./config.json /app
COPY ./shadowsocks-local /app
COPY ./source.list /app

RUN mv source.list /etc/apt/


RUN apt-get update && apt-get upgrade && apt-get install polipo && \
    service  polipo  start && \
    polipo socksParentProxy=127.0.0.1:1080  && \
    service  polipo stop  && \
    service  polipo  start

ENTRYPOINT ["./shadowsocks-local"]

EXPOSE 8123
EXPOSE 1080
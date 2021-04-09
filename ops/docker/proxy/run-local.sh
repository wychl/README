#!/bin/bash

docker run --name proxy --restart always --network=host -v ${PWD}/local.json:/etc/shadowsocks.json -p 1081:1080 -d proxy-local:v1
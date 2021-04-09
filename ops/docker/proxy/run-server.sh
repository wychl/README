#!/bin/bash

docker run --name proxy --restart always --network=host -v /path/to/config.json:/etc/shadowsocks.json -p local_port:container_port -d proxy-server:v1
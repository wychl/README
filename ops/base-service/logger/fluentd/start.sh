#!/bin/sh

# 启动fluentd服务
docker run --name fluentd -d -p 9880:9880 -v ${PWD}/conf/:/fluentd/etc/ -v ${PWD}/data:/fluentd/log fluent/fluentd:v1.3-debian-1

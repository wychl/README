# fluentd

## 启动fluentd实例

- 配置文件 `conf/fluent.conf`
  - 输入插件`http`,端口9880
  - 输出插件 `file` match为`sample.*`,日志文件位置：`/fluentd/log/sample`
- 数据文件在 `data/`
- 修改`data/`文件权限  `chmod 777 data/`



```sh
docker run --name fluentd -d -p 9880:9880 -v ${PWD}/conf/:/fluentd/etc/ -v ${PWD}/data:/fluentd/log  fluent/fluentd:v1.3-debian-1
```

- 测试

```
curl -X POST -d 'json={"json":"message"}' http://localhost:9880/sample.test
```


## 记录docker容器日志

- 启动fluentd ` docker run --name fluentd -d -p 24224:24224 -v ${PWD}/conf/:/fluentd/etc/ -v ${P
WD}/data:/fluentd/log  fluent/fluentd:v1.3-debian-1`
- 测试 `docker run  --rm --log-driver=fluentd --log-opt tag="docker.{{.ID}}" alpine echo 'Hello Fluentd!'`


## 记录kubernetes pod日志

**注意fluentd以sidecar方式和应用部署在同一个pod中**

```sh
kubectl create -f k8s/config.yaml
kubectl crate -f k8s/demo.yaml
```

## golang应用日志

- 启动fluentd ` docker run --name fluentd -d -p 24224:24224 -v ${PWD}/conf/:/fluentd/etc/ -v ${P
WD}/data:/fluentd/log  fluent/fluentd:v1.3-debian-1`
- 测试 `go run main.go`
## istio介绍


### Envoy(代理)

Envoy 作为为sidecar和对应服务部署在在同一个pod中。


```sh
docker ps | grep strprocesser-v1

    41a606c6a25b        8d28edd1b0ff                   "/usr/local/bin/pilo…"   2 hours ago         Up 2 hours                              k8s_istio-proxy_strprocesser-v1-5b6f86767f-2s7b6_default_cdbdcde4-5549-11e9-861e-00163e02a016_0
    d385ae5c7baf        wanyanchengli/strprocesser     "./server"               2 hours ago         Up 2 hours                              k8s_strprocesser_strprocesser-v1-5b6f86767f-2s7b6_default_cdbdcde4-5549-11e9-861e-00163e02a016_0
    ab9fb2e88a43        k8s.gcr.io/pause:3.1           "/pause"                 2 hours ago         Up 2 hours                              k8s_POD_strprocesser-v1-5b6f86767f-2s7b6_default_cdbdcde4-5549-11e9-861e-00163e02a016_0
```

## 网络配置

```sh
kubectl --insecure-skip-tls-verify apply -f networking/destination.yaml
kubectl --insecure-skip-tls-verify apply -f networking/virtualservice-all.yaml
kubectl --insecure-skip-tls-verify  apply -f networking/gateway.yaml
```

## 白名单
```sh
kubectl --insecure-skip-tls-verify apply -f policy/whitelist.yaml
kubectl --insecure-skip-tls-verify delete -f policy/whitelist.yaml
```

## 日志(Fluentd)

```sh
#stdio

kubectl --insecure-skip-tls-verify apply -f logging/new_logs.yaml
kubectl --insecure-skip-tls-verify logs -n istio-system -l istio-mixer-type=telemetry -c mixer | grep \"instance\":\"newlog.logentry.istio-system\" | grep -v '"destination":"telemetry"' | grep -v   '"destination":"pilot"' | grep -v '"destination":"policy"' | grep -v '"destination":"unknown"'
kubectl --insecure-skip-tls-verify delete -f logging/new_logs.yaml

## fluentd
kubectl --insecure-skip-tls-verify apply -f logging/logging.yaml
kubectl --insecure-skip-tls-verify apply -f logging/fluentd-istio.yaml
kubectl  --insecure-skip-tls-verify -n logging port-forward $(kubectl --insecure-skip-tls-verify -n logging get pod -l app=kibana -o jsonpath='{.items[0].metadata.name}') 5601:5601 &
kubectl --insecure-skip-tls-verify delete -f logging/fluentd-istio.yam


```

## metric

```sh
kubectl --insecure-skip-tls-verify apply -f metric/metric.yaml

kubectl --insecure-skip-tls-verify -n istio-system port-forward $(kubectl --insecure-skip-tls-verify -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 9090:9090 &

kubectl --insecure-skip-tls-verify delete -f metric/metric.yaml

```

`prometheus`查询`istio_double_request_count`

## 通过grafana查看metric

```sh
kubectl  --insecure-skip-tls-verify -n istio-system port-forward $(kubectl  --insecure-skip-tls-verify -n istio-system get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000 &
```


地址  http://localhost:3000/d/fZeFu16mz/istio-mesh-dashboard?




## 分布式追踪

```sh
#jaeger
kubectl --insecure-skip-tls-verify  port-forward -n istio-system $(kubectl --insecure-skip-tls-verify  get pod -n istio-system -l app=jaeger -o jsonpath='{.items[0].metadata.name}') 16686:16686  &
```

地址 http://localhost:16686/

## 断路器

```sh
kubectl --insecure-skip-tls-verify apply -f https://raw.githubusercontent.com/istio/istio/release-1.1/samples/httpbin/sample-client/fortio-deploy.yaml
FORTIO_POD=$(kubectl --insecure-skip-tls-verify get pod | grep fortio | awk '{ print $1 }')
kubectl --insecure-skip-tls-verify exec -it ${FORTIO_POD}  -c fortio /usr/bin/fortio -- load -curl  http://strprocesserui:9090/ui #没有触发

kubectl --insecure-skip-tls-verify exec -it $FORTIO_POD  -c fortio /usr/bin/fortio -- load -c 2 -qps 0 -n 80 -loglevel Warning  http://strprocesserui:9090/ui #触发


```


## 问题

Quota excluded. service: strprocesserui.default.svc.cluster.local matches binding shortname: strprocesser, but does not match fqdn: strprocesser.default.svc.cluster.local
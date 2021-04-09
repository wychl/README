# prometheus poc

Prometheus是一个开源系统监控和警报工具包，

## 特点

- 具有由度量名称和键/值对标识的时间序列数据的多维数据模型
- PromQL 查询语言
- 不依赖分布式存储和其他的服务;单个服务器节点是自治的
- 通过HTTP轮询方式采集时间序列数据
- 通过中间网关推送时间序列数据
- 通过服务发现和静态配置采集目标
- 快速诊断问题，提高系统的可靠性
- 当配置文件改变时，实时加载

## 组件

- server 采集和存储数据
- client 度量采集目标指标
- a push gateway for supporting short-lived jobs
- special-purpose exporters for services like HAProxy, StatsD, Graphite, etc.
- alertmanager 警报系统
- various support tools

## 适合场景

- 记录任何数字时间序列
- 适合机器为中心的监控
- 微服务架构
- 它对多维数据收集和查询是一种特殊的优势

## 不适合场景

- 需要100%的准确度，比如：按请求次数计费场景（promethemus收集的数据既不完整和也不详细）


## 安装prometheus

- 启动命令

```sh
docker run \
    --name prometheus \
    -p 9090:9090 \
    -v ${PWD}/prometheus:/etc/prometheus \
    -d \
    prom/prometheus:v2.9.2 \
    --config.file=/etc/prometheus/prometheus.yml
```

- UI地址 http://127.0.0.1:9090/graph

## 安装influxdb

```sh
docker run \
    --name influxdb \
    -v $PWD/influxdb-data:/var/lib/influxdb \
    -v $PWD/influxdb-config:/etc/influxdb \
    -e INFLUXDB_DB=prometheus \
    -e INFLUXDB_ADMIN_ENABLED=true \
    -e INFLUXDB_ADMIN_USER=admin \
    -e INFLUXDB_ADMIN_PASSWORD=admin \
    -e INFLUXDB_USER=prometheus \
    -e INFLUXDB_USER_PASSWORD=prometheus \
    -e INFLUXDB_GRAPHITE_ENABLED=true \
    -p 8086:8086 \
    -p 8083:8083 \
    -p 2003:2003 \
    -d \
    influxdb:1.7.6
```



### 生成`influxdb`默认配置

```sh
docker run --rm influxdb:1.7.6 influxd config > influxdb-config/influxdb.conf
```

## `grafana`安装和配置

### 启动grafana

```sh
 docker run \
    -v $PWD/grafana-data:/var/lib/grafana \
    -d \
    -p 3000:3000 \
    --name grafana \
    -e "GF_SECURITY_ADMIN_PASSWORD=admin" \
    grafana/grafana:5.1.0
```

### 配置`influxdb`数据源 

## 使用`表达式`查看prometheus本身的监控日志

- `promhttp_metric_handler_requests_total`
    `/metrics`路由请求数量
- `promhttp_metric_handler_requests_total{code="200"}`
    `/metrics`路由请求数量（http code为200的）
- `count(promhttp_metric_handler_requests_total)`
    产生的时间序列数量

## 启动 alterManager

```sh
docker run \
    --name alertmanager \
    -v ${PWD}/alertmanager:/etc/alertmanager \
    -p 9093:9093 \
    -p 9094:9094 \
    -p 9095:9095 \
    -d \
    quay.io/prometheus/alertmanager:v0.17.0 \
    --config.file=/etc/alertmanager/alertmanager.yml

```

## 设置警报规则（rules.yml文件）

事例1

```yaml
groups:
- name: example
  rules:
  - alert: HighErrorRate
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels: #警报附加标签
      severity: page
    annotations: #更长的附加信息
      summary: High request latency
```

模版配置事例

```yaml
groups:
- name: example
  rules:

  # Alert for any instance that is unreachable for >5 minutes.
  - alert: InstanceDown
    expr: up == 0
    for: 5m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} down"
      description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."

  # Alert for any instance that has a median request latency >1s.
  - alert: APIHighRequestLatency
    expr: api_http_request_latencies_second{quantile="0.5"} > 1
    for: 10m
    annotations:
      summary: "High request latency on {{ $labels.instance }}"
      description: "{{ $labels.instance }} has a median request latency above 1s (current value: {{ $value }}s)"
```

## InfluxDB作为Prometheus的后端存储

### InfluxDB api

- `/api/v1/prom/read`
- `/api/v1/prom/write`
- `/api/v1/prom/metrics`

### Create a target database

Create a database in your InfluxDB instance to house data sent from Prometheus. In the examples provided below, prometheus is used as the database name, but you’re welcome to use the whatever database name you like.

`CREATE DATABASE "prometheus"`

### Configuration

To enable the use of the Prometheus remote read and write APIs with InfluxDB, add URL values to the following settings in the Prometheus configuration file:

- simple

```yaml
remote_write:
  - url: "http://localhost:8086/api/v1/prom/write?db=prometheus"

remote_read:
  - url: "http://localhost:8086/api/v1/prom/read?db=prometheus"
```

- authentication

```yaml
remote_write:
  - url: "http://localhost:8086/api/v1/prom/write?db=prometheus&u=username&p=password"
remote_read:
  - url: "http://localhost:8086/api/v1/prom/read?db=prometheus&u=username&p=password"
```

## 使用prometheus和grafana监控系统

### 安装 Node Exporter

```sh
docker run \
  --name node-exporter \
  -d -p 9100:9100 \
  -v "/proc:/host/proc" \
  -v "/sys:/host/sys" \
  -v "/:/rootfs" \
  prom/node-exporter
```

### 添加 prometheus data source

1. Name: Prometheus
2. Type: Prometheus
3. URL: http://:9090, (default port is 9090)
4. Access: proxy
5. Basic Auth: According to Your Server

### 导入Prometheus Stats仪表板

1. 下载仪表盘配置 curl -o prometheus-dash.json https://grafana.com/api/dashboards/3662/revisions/2/download
2. 将配置导入到 grafana

## 表达式

### 将抓取的数据聚合成新的时间序列

- 表达式 `avg(rate(rpc_durations_seconds_count[5m])) by (job, service)`

## alertmanager

### 抑制机制

Alertmanager的抑制机制可以避免当某种问题告警产生之后用户接收到大量由此问题导致的一系列的其它告警通知。例如当集群不可用时，用户可能只希望接收到一条告警，告诉他这时候集群出现了问题，而不是大量的如集群中的应用异常、中间件服务异常的告警通知。

例如，定义如下抑制规则：

```yaml
- source_match:
    alertname: NodeDown
    severity: critical
  target_match:
    severity: critical
  equal:
    - node
```

例如当集群中的某一个主机节点异常宕机导致告警NodeDown被触发，同时在告警规则中定义了告警级别severity=critical。由于主机异常宕机，该主机上部署的所有服务，中间件会不可用并触发报警。根据抑制规则的定义，如果有新的告警级别为severity=critical，并且告警中标签node的值与NodeDown告警的相同，则说明新的告警是由NodeDown导致的，则启动抑制机制停止向接收器发送通知。

概念介绍

## Mixer提供三个核心功能：

1. 前置条件检查（Precondition Checking）：某一服务响应外部请求前，通过Envoy向Mixer发送Check请求，检查该请求是否满足一定的前提条件，包括白名单检查、ACL检查等。
2. 配额管理（Quota Management）：当多个请求发生资源竞争时，通过配额管理机制可以实现对资源的有效管理。
3. 遥测报告上报（Telemetry Reporting）：该服务处理完请求后，通过Envoy向Mixer上报日志、监控等数据。

## Attribute（属性）

大部分attributes由Envoy提供。Istio用attributes来控制服务在Service Mesh中运行时行为。attributes是有名称和类型的元数据，用来描述入口和出口流量和流量产生时的环境。attributes携带了一些具体信息，比如：API请求状态码、请求响应时间、TCP连接的原始地址等。

## RefrencedAttributes（被引用的属性）

refrencedAttributes是Mixer Check时进行条件匹配后被使用的属性的集合。Envoy向Mixer发送的Check请求中传递的是属性的全集，refrencedAttributes只是该全集中被应用的一个子集。

举个例子，Envoy某次发送的Check请求中发送的attributes为{request.path: xyz/abc, request.size: 234,source.ip: 192.168.0.1}，如Mixer中调度到的多个adapters只用到了request.path和request.size这两个属性。那么Check后返回的refrencedAttributes为{request.path: xyz/abc, request.size: 234}。
为防止每次请求时Envoy都向Mixer中发送Check请求，Mixer中建立了一套复杂的缓存机制，使得大部分请求不需要向Mixer发送Check请求。

request.path: xyz/abc
request.size: 234
request.time: 12:34:56.789 04/17/2017
source.ip: 192.168.0.1
destination.service: example

属性词汇由[_.a-z0-9]组成，其中"."为命名空间分隔符，所有属性词汇可以查看这里，属性类型可以查看这里。

## Adapter（适配器）

Mixer是一个高度模块化、可扩展组件，内部提供了多个适配器(adapter)。
Envoy提供request级别的属性（attributes）数据。
adapters基于这些attributes来实现日志记录、监控指标采集展示、配额管理、ACL检查等功能。Istio内置的部分adapters举例如下：

## Template（模板）

对于一个网络请求，Mixer通常会调用两个rpc：Check和Report。不同的adapter需要不同的attributes，template定义了attributes到adapter输入数据映射的schema，一个适配器可以支持多个template。一个上报metric数据的模板如下所示：

```yaml
apiVersion: "config.istio.io/v1alpha2"
kind: metric
metadata:
  name: requestsize
  namespace: istio-system
spec:
  value: request.size | 0
  dimensions:
    source_service: source.service | "unknown"
    source_version: source.labels["version"] | "unknown"
    destination_service: destination.service | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    response_code: response.code | 200
  monitored_resource_type: '"UNSPECIFIED"'
```

模板字段的值可以是字面量或者表达式，如果时表达式，则表达式的值类型必须与字段的数据类型一致。

**Mixer的配置模型**

Mixer的yaml配置可以抽象成三种模型：Handler、Instance、Rule

这三种模型主要通过yaml中的kind字段做区分，kind值有如下几种：

- adapter kind：表示此配置为Handler。
- template kind：表示此配置为Template。
- "rule"：表示此配置为Rule。

## Handler

一个Handler是配置好的Adpater的实例。Handler从yaml配置文件中取出adapter需要的配置数据。一个典型的Promethues Handler配置如下所示：

```yaml
apiVersion: config.istio.io/v1alpha2
kind: prometheus
metadata:
  name: handler
  namespace: istio-system
spec:
  metrics:
  - name: request_count
    instance_name: requestcount.metric.istio-system
    kind: COUNTER
    label_names:
    - destination_service
    - destination_version
    - response_code
```


## Instance

Instance定义了attributes到adapter输入的映射，一个处理requestduration metric数据的Instance配置如下所示：

```yaml
apiVersion: config.istio.io/v1alpha2
kind: metric
metadata:
  name: requestduration
  namespace: istio-system
spec:
  value: response.duration | "0ms"
  dimensions:
    destination_service: destination.service | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    response_code: response.code | 200
  monitored_resource_type: '"UNSPECIFIED"'
```

上述Instance的完全限定名是requestduration.metric.istio-system，Handler和Rule可以通过这个名称对此Instance进行引用。

## Rule

Rule定义了一个特定的Instance何时调用一个特定的Handler，一个典型的Rule配置如下所示：

```yaml
apiVersion: config.istio.io/v1alpha2
kind: rule
metadata:
  name: promhttp
  namespace: istio-system
spec:
  match: destination.service == "service1.ns.svc.cluster.local" && request.headers["x-user"] == "user1"
  actions:
  - handler: handler.prometheus
    instances:
    - requestduration.metric.istio-system
```

上述例子中，定义的Rule为：对目标服务为service1.ns.svc.cluster.local且request.headers["x-user"] 为user1的请求，Instance: requestduration.metric.istio-system才调用Handler: handler.prometheus。

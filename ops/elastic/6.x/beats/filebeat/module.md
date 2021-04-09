
# module

## 概述

- Filebeat提供了一组预先构建的模块，你可以使用这些模块快速实现并部署一个日志监控解决方案，包括样例指示板和数据可视化
- 默认路径 /etc/filebeat/modules.d/

## 参考链接

- https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-modules.html

## redis

### 配置

```yaml
- module: redis
  log:
    enabled: true
    var.paths: ["/path/to/log/redis/redis-server.log*"]
  slowlog:
    enabled: true
    var.hosts: ["localhost:6378"]
    var.password: "YOUR_PASSWORD"
```
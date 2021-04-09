# MySQL日志

## 开启慢查询日志

编辑 `/etc/my.cnf`文件

```conf
# vim /etc/my.cnf
[mysqld]
log_queries_not_using_indexes = 1
slow_query_log=on
slow_query_log_file=/var/log/mysql/slow-mysql-query.log
long_query_time=0

[mysqld_safe]
log-error=/var/log/mysql/mysqld.log
```

## 配置filebeat

### 配置 Filebeat modules 动态加载。

```yaml
filebeat.config.modules:
  # Glob pattern for configuration loading
  path: /etc/filebeat/modules.d/mysql.yml

  # Set to true to enable config reloading
  reload.enabled: true

  # Period on which files under path should be checked for changes
  reload.period: 1s
```

## 填入Kibana地址

```yaml
setup.kibana:

  # Kibana Host
  # Scheme and port can be left out and will be set to the default (http and 5601)
  # In case you specify and additional path, the scheme is required: http://localhost:5601/path
  # IPv6 addresses should always be defined as: https://[2001:db8::1]:5601
  host: "https://es-cn-0p11111000zvqku.kibana.elasticsearch.aliyuncs.com:5601"
```


## 填入 Elasticsearch 地址及端口号

```yaml
output.elasticsearch:
  # Array of hosts to connect to.
  hosts: ["es-cn-0p11111000zvqku.elasticsearch.aliyuncs.com:9200"]
  # Optional protocol and basic auth credentials.
  #protocol: "https"
  username: "elastic"
  password: "elastic@333"
```

## 启用 MySQL 模块，并进行配置：


```yaml
# sudo filebeat modules enable mysql
# vim /etc/filebeat/modules.d/mysql.yml
- module: mysql
  # Error logs
  error:
    enabled: true
    var.paths: ["/var/log/mysql/mysqld.log"]
    # Set custom paths for the log files. If left empty,
    # Filebeat will choose the paths depending on your OS.
    #var.paths:

  # Slow logs
  slowlog:
    enabled: true
    var.paths: ["/var/log/mysql/slow-mysql-query.log"]
    # Set custom paths for the log files. If left empty,
    # Filebeat will choose the paths depending on your OS.
    #var.paths:
```
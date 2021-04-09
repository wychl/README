# 使用

## 命令列表

- https://www.elastic.co/guide/en/beats/filebeat/6.8/command-line-options.html

| 命令     | 描述                                                                           |
| -------- | ------------------------------------------------------------------------------ |
| export   | Exports the configuration, index template, or a dashboard to stdout.           |
| help     | 帮助命令                                                                       |
| keystore | secrets管理命令                                                                |
| modules  | modules管理命令                                                                |
| run      | 启动Filebeat；默认执行的命令                                                   |
| setup    | 设置初始化环境变量（index template, Kibana dashboards ,machine learning jobs ( |
| test     | 测试配置 `filebeat test config -c filebeat.yml`                                |
| version  | 查询版本                                                                       |

## modules

### 启动modules三种方式

#### 在`modules.d`目录中启动

```sh
filebeat modules enable apache2 mysql

# 查看列表
filebeat modules list
```

#### `filebeat`启动时指定

```sh
filebeat --modules nginx,mysql,system
```

#### 在`filebeat.yml`配置文件中

```yaml
filebeat.modules:
- module: nginx
- module: mysql
- module: system
```

## 配置

### 输入

#### 参考链接

- https://www.elastic.co/guide/en/beats/filebeat/6.8/configuration-filebeat-options.html

#### 不同输入类型的配置

- [Log](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-log.html)
- [Stdin](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-stdin.html)
- [Redis](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-redis.html)
- [UDP](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-udp.html)
- [Docker](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-docker.html)
- [TCP](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-tcp.html)
- [Syslog](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-syslog.html)
- [NetFlow](https://www.elastic.co/guide/en/beats/filebeat/6.8/filebeat-input-netflow.html)

#### 通用输入配置

链接 https://www.elastic.co/guide/en/beats/filebeat/6.8/configuration-general-options.html

- `name` 
  - Beat名称
  - 空值时，使用服务器的`hostname`
  - 聚合
- `tags`
  - 聚合
- `fields`
  - `fields_under_root`
  - 添加字段到日志数据中
- `processors`
  - 处理方法

### 多行转为一行

- https://www.elastic.co/guide/en/beats/filebeat/6.8/multiline-examples.html

场景：

- Java stack trace

```plaintext
Exception in thread "main" java.lang.IllegalStateException: A book has a null property
       at com.example.myproject.Author.getBookIds(Author.java:38)
       at com.example.myproject.Bootstrap.main(Bootstrap.java:14)
Caused by: java.lang.NullPointerException
       at com.example.myproject.Book.getId(Book.java:22)
       at com.example.myproject.Author.getBookIds(Author.java:35)
       ... 1 more
```

- C-style line

```plaintext
printf ("%10.10ld  \t %10.10ld \t %s\
  %f", w, x, y, z );
```

- 时间戳日志

```plaintext
[2015-08-24 11:49:14,389][INFO ][env                      ] [Letha] using [1] data paths, mounts [[/
(/dev/disk1)]], net usable_space [34.5gb], net total_space [118.9gb], types [hfs]
```

### 配置队列

#### 内存队列

内存队列将接收到的日志事件首先保留在内存缓冲区中，当满足设置条件时再发送出去。

```plaintext
queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 5s
```

- `events`
  - 默认值：4096
  - 可以存储的日志事件数
- `flush.min_events`
  - 默认值： 2048
  - 发出所需的最小日志事件数
  - 0值时，进来的日志立即发出
- `flush.timeout`
  - 当接收日志事件数量已经达到`flush.min_events`设置的值，要发出还需要等待的时间
  - 值为0时，不等待立即发出

#### spool queue

- https://www.elastic.co/guide/en/beats/filebeat/6.8/configuring-internal-queue.html#configuration-internal-queue-spool

```plaintex
queue.spool:
  file:
    path: "${path.data}/spool.dat"
    size: 512MiB
    page_size: 16KiB
  write:
    buffer_size: 10MiB
    flush.timeout: 5s
    flush.events: 1024
```

- `file.path`
  - 默认值 `"${path.data}/spool.dat".`
  - `spool`文件路径
  - 没有时，启动时创建
- `file.permissions`
  - 文件权限
  - 默认值:`0600`
- [`file.size`](https://www.elastic.co/guide/en/beats/filebeat/6.8/configuring-internal-queue.html#_file_size)
  - 文件大小
  - 默认值: 100MB
  - 仅在文件创建时设置，且不能更改
- `file.prealloc`
  - 默认值: true
  - 仅在文件创建时设置，且不能更改
  - 如果prealloc设置为false，则文件将动态增长。如果prealloc为false并且系统磁盘空间不足，则假脱机将阻塞
- `file.buffer.size`
  - 写入缓冲区大小
  - 一旦超出缓冲区大小，将flush write buffer
  - 默认：1MB
- `write.codec`
  - 日志事件序列化方式
  - 支持json和cbor
  - 默认值：cbor
- `write.flush.timeout`
  - flush操作的等待时间
  - 默认值： 1s
  - 值为0时，立即flush
- `write_flush_events`
  - 日志事件数量，超过的立即flush
  - 默认值： 16384
- `read.flush.timeout`
  - 默认值:0s
  - 值为0s,日志立即发出
# 备份与恢复

要使用这个功能，你必须首先创建一个保存数据的仓库。有多个仓库类型可以供你选择：

- 共享文件系统，比如 NAS
- Amazon S3 阿里云OSS
- HDFS (Hadoop 分布式文件系统)
- Azure Cloud

## 参考链接

- https://www.elastic.co/guide/cn/elasticsearch/guide/current/backing-up-your-cluster.html
- https://www.elastic.co/guide/cn/elasticsearch/guide/current/_restoring_from_a_snapshot.html
- https://developer.aliyun.com/article/708959
- https://www.elastic.co/guide/en/elasticsearch/plugins/6.8/repository-s3.html?spm=a2c6h.12873639.0.0.dc4c24e9RkDUVq
- https://www.elastic.co/guide/en/elasticsearch/plugins/6.8/repository.html

## 备份一般过程

### 创建仓库

```sh
curl -X PUT "localhost:9200/_snapshot/my_backup" -H 'Content-Type: application/json' -d'
{
    "type": "fs", 
    "settings": {
        "location": "/mount/backups/my_backup" 
         "max_snapshot_bytes_per_sec" : "20mb", 
        "max_restore_bytes_per_sec" : "20mb"
    }
}
'
```

- `my_backup` 备份仓库取名字，在本例它叫 my_backup 。
- `type` 我们指定仓库的类型应该是一个共享文件系统。
- `location` 已挂载的设备作为目的地址
- max_snapshot_bytes_per_sec
当快照数据进入仓库时，这个参数控制这个过程的限流情况。默认是每秒 20mb
- max_restore_bytes_per_sec
当从仓库恢复数据时，这个参数控制什么时候恢复过程会被限流以保障你的网络不会被占满。默认是每秒 20mb。


***注意：共享文件系统路径必须确保集群所有节点都可以访问到。***

### 快照所有索引

```sh
curl -X PUT "localhost:9200/_snapshot/my_backup/snapshot_1" \
  -H 'Content-Type: application/json'
```

这个会备份所有打开的索引到 my_backup 仓库下一个命名为 snapshot_1 的快照里。这个调用会立刻返回，然后快照会在后台运行。

### 快照指定索引

```sh
curl -X PUT "localhost:9200/_snapshot/my_backup/snapshot_2" \
  -H 'Content-Type: application/json' \
  -d'
    {
        "indices": "index_1,index_2"
    }'
```

### 列出快照信息

```sh
curl -X PUT "localhost:9200/_snapshot/my_backup/snapshot_2" \
  -H 'Content-Type: application/json'
```

### 删除快照

```sh
curl -X DELETE "localhost:9200/_snapshot/my_backup/snapshot_2" \
  -H 'Content-Type: application/json'
```

## 恢复

### 恢复快照里的所有索引

```sh
curl -X POST "localhost:9200/_snapshot/my_backup/snapshot_1/_restore" \
  -H 'Content-Type: application/json'
```

### 恢复指定的索引


```sh
curl -X POST "localhost:9200/_snapshot/my_backup/snapshot_1/_restore" \
  -H 'Content-Type: application/json' \
  -d '{
    "indices": "index_1", 
    "rename_pattern": "index_(.+)", 
    "rename_replacement": "restored_index_$1" 
}'
```

***这个会恢复 index_1 到你及群里，但是重命名成了 restored_index_1 。***

- 只恢复 `index_1` 索引，忽略快照中存在的其余索引。
- 查找所提供的模式能匹配上的正在恢复的索引。
- 然后把它们重命名成替代的模式。
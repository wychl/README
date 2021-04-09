
# 索引

## 参考文档

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices.html
- https://www.elastic.co/guide/cn/elasticsearch/guide/current/index-management.html

## 索引名称限制

- 小写字母
- 不能包含 “ \, /, *, ?, ", <, >, |, ` ` (space character), ,, #” 字符
- `:`字符在7.0不再支持
- 不能以` -, _, +`字符开始
- Cannot be . or ..
- 不能超过`255`个字节

## 创建

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-create-index.html

### API

```sh
curl -X PUT "localhost:9200/twitter?pretty" -H 'Content-Type: application/json' -d'
{
    "settings" : {
        "index" : {
            "number_of_shards" : 3, 
            "number_of_replicas" : 2 
        }
    }
}
'
```

- `number_of_shards`
每个索引的主分片数，默认值是 5 。这个配置在索引创建后不能修改。
- `number_of_replicas`
每个主分片的副本数，默认值是 1 。对于活动的索引库，这个配置可以随时修改。

### 修改主分片副本数

```sh
curl -X PUT "localhost:9200/twitter/_settings?pretty" -H 'Content-Type: application/json' -d'
{
    "number_of_replicas": 1
}
'
```

## 更新

## 查看

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-get-index.html

### API

```sh
# 支持索引名称，别名，通配符
curl -X GET "localhost:9200/twitter?pretty"

# 检索所有的索引
curl -X GET "localhost:9200/_all?pretty"
curl -X GET "localhost:9200/*?pretty"
```

## 检测索引是否存在

### API

```sh
# 检测索引是否存在
curl -I "localhost:9200/twitter?pretty"
```

***跟胡HTTP相映的状态代码指示索引是否存在。 404表示它不存在，而200表示它不存在。***

## 索引开启和关闭

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-open-close.html

### API

```sh
# 操作单个
curl -X POST "localhost:9200/my_index/_close?pretty"
curl -X POST "localhost:9200/my_index/_open?pretty"

# 操作所有的索引
curl -X POST "localhost:9200/_all/_close?pretty"
curl -X POST "localhost:9200/*/_open?pretty"
```

- ***为了降低数据丢失的风险，请避免将封闭的索引长时间保留在群集中，因为当节点离开集群或者被替换时关闭的索引不会再被创建。***

### 禁止关闭索引

在`elasticsearch.yml`文件中做如下配置：

```yaml
cluster.indices.close.enable: false
```

## 删除

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-delete-index.html
- https://www.elastic.co/guide/cn/elasticsearch/guide/current/_deleting_an_index.html

### API

```sh
# 删除单个
curl -X DELETE "localhost:9200/twitter?pretty"


# 删除多个
curl -X DELETE "localhost:9200/index_one,index_two?pretty"
curl -X DELETE "localhost:9200/index_*?pretty"

# 删除所有
curl -X DELETE "localhost:9200/_all?pretty"
curl -X DELETE "localhost:9200/*?pretty"
```

## 禁止通过通配符或_all标识索引

在`elasticsearch.yml`文件中做如下配置：

```yaml
action.destructive_requires_name: true
```

## shrink index 

目的：shrink an existing index into a new index with fewer primary shards.

***The requested number of primary shards in the target index must be a factor of the number of shards in the source index. For example an index with 8 primary shards can be shrunk into 4, 2 or 1 primary shards or an index with 15 primary shards can be shrunk into 5, 3 or 1.***

## 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-shrink-index.html

### 原理

- First, it creates a new target index with the same definition as the source index, but with a smaller number of primary shards.
- Then it hard-links segments from the source index into the target index. (If the file system doesn’t support hard-linking, then all segments are copied into the new index, which is a much more time consuming process. Also if using multiple data paths, shards on different data paths require a full copy of segment files if they are not on the same disk since hardlinks don’t work across disks)
- Finally, it recovers the target index as though it were a closed index which had just been re-opened.

### _preparing_an_index_for_shrinking

```sh
curl -X PUT "localhost:9200/my_source_index/_settings?pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "index.routing.allocation.require._name": "shrink_node_name", 
    "index.blocks.write": true 
  }
}
'
```

- `index.routing.allocation.require._name` 强制将每个分片的副本重定位到名称为`shrink_node_name`的节点
- `index.blocks.write`，阻止对索引的写操作

### _shrinking_an_index

#### 要求

The target index must not exist.
The source index must have more primary shards than the target index.
The number of primary shards in the target index must be a factor of the number of primary shards in the source index. The source index must have more primary shards than the target index.
The index must not contain more than 2,147,483,519 documents in total across all shards that will be shrunk into a single shard on the target index as this is the maximum number of docs that can fit into a single shard.
The node handling the shrink process must have sufficient free disk space to accommodate a second copy of the existing index.

#### API

```sh
curl -X POST "localhost:9200/my_source_index/_shrink/my_target_index?copy_settings=true&pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "index.routing.allocation.require._name": null, 
    "index.blocks.write": null 
  }
}
'

curl -X POST "localhost:9200/my_source_index/_shrink/my_target_index?copy_settings=true&pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "index.number_of_replicas": 1,
    "index.number_of_shards": 1, 
    "index.codec": "best_compression" 
  },
  "aliases": {
    "my_search_indices": {}
  }
}
'

```

## 拆分索引

拆分索引允许您将现有索引拆分为新索引，在该索引中，每个原始主碎片在新索引中均拆分为两个或多个主碎片。

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-split-index.html

## PUT mapping

允许您将字段添加到现有索引或更改现有字段的搜索设置

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-put-mapping.html

## 索引别名

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/indices-aliases.html

## 索引生命周期

### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/index-lifecycle-management.html


### 生命周期策略功能

- The maximum size or age at which you want to roll over to a new index.
- The point at which the index is no longer being updated and the number of primary shards can be reduced.
- When to force a merge to permanently delete documents marked for deletion.
- The point at which the index can be moved to less performant hardware.
- The point at which the availability is not as critical and the number of replicas can be reduced.
- When the index can be safely deleted.

### 示例

#### 创建策略

```sh
curl -X PUT "localhost:9200/_ilm/policy/datastream_policy?pretty" -H 'Content-Type: application/json' -d'
{
  "policy": {                       
    "phases": {
      "hot": {                      
        "actions": {
          "rollover": {             
            "max_size": "50GB",
            "max_age": "30d"
          }
        }
      },
      "delete": {
        "min_age": "90d",           
        "actions": {
          "delete": {}              
        }
      }
    }
  }
}
'
```

- `datastream_policy` 策略名称
- `phases.hot`  hot阶段定义； `min_age`可选字段，默认值0ms
- `rollover` 索引写入的达到50GB或使用30天后对其进行`rollover`
- `phases.delete` 删除阶段90天之后开始
- `actions.delete` 删除动作定义

#### 策略与索引关联

- 第一种方式：通过索引模版

```sh
curl -X PUT "localhost:9200/_template/datastream_template?pretty" -H 'Content-Type: application/json' -d'
{
  "index_patterns": ["datastream-*"],                 
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "index.lifecycle.name": "datastream_policy",      
    "index.lifecycle.rollover_alias": "datastream"    
  }
}
'
```

- `index_patterns`，索引匹配方式；示例中匹配所有以`datastream-`开头的索引
- `index.lifecycle.name` 策略名称
- `index.lifecycle.rollover_alias`,`rollover action`别名

- 第二种方式 创建索引时

```sh
curl -X PUT "localhost:9200/test-index?pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "index.lifecycle.name": "my_policy"
  }
}
'
```

#### 创建可以匹配策略的测试索引

```sh
curl -X PUT "localhost:9200/datastream-000001?pretty" -H 'Content-Type: application/json' -d'
{
  "aliases": {
    "datastream": {
      "is_write_index": true
    }
  }
}
'
```

- 策略中使用了`Rollover Action`操作，索引索引名称以数字结尾(`000001`)
- `aliases` 索引别名，与上步骤中的`index.lifecycle.rollover_alias`字段一致

#### 查看索引生命周期信息

```sh
curl -X GET "localhost:9200/datastream-*/_ilm/explain?pretty"
```

返回值

```json
{
  "indices": {
    "datastream-000001": {
      "index": "datastream-000001",
      "managed": true,                           
      "policy": "datastream_policy",             
      "lifecycle_date_millis": 1538475653281,
      "phase": "hot",                            
      "phase_time_millis": 1538475653317,
      "action": "rollover",                      
      "action_time_millis": 1538475653317,
      "step": "attempt-rollover",                
      "step_time_millis": 1538475653317,
      "phase_execution": {
        "policy": "datastream_policy",
        "phase_definition": {                    
          "min_age": "0ms",
          "actions": {
            "rollover": {
              "max_size": "50gb",
              "max_age": "30d"
            }
          }
        },
        "version": 1,                            
        "modified_date_in_millis": 1539609701576
      }
    }
  }
}
```

### 策略阶段和动作

#### 参考链接

- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html
- https://www.elastic.co/cn/blog/implementing-hot-warm-cold-in-elasticsearch-with-index-lifecycle-management

#### 阶段

- `hot` 积极更新和查询
  - [Set Priority](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-set-priority-action)
  - [Rollover](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-rollover-action)
  - [Unfollow](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-unfollow-action)
- `warm` 不再更新但是仍在查询
  - [Set Priority](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-set-priority-action)
  - [Allocate](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-allocate-action)
  - [Read-Only](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-readonly-action)
  - [Force Merge](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-forcemerge-action)
  - [Shrink](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-shrink-action)
  - [Unfollow](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-unfollow-action)
- `cold` 再更新但是有少量的查询
  - [Set Priority](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-set-priority-action)
  - [Allocate](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-allocate-action)
  - [Freeze](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-freeze-action)
  - [Unfollow](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-unfollow-action)
- `delete` 不再需要，可以安全的删除
  - [Delete](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/_actions.html#ilm-delete-action)

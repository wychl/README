# filebeat

## 装备工作

### 采用basic auth的方式链接elasticsearch

#### 参考链接

- https://www.elastic.co/guide/en/beats/filebeat/6.8/beats-basic-auth.html

#### 创建链接账号和配置权限

- ***行环境：类Unix系统***
- ***可以在任何目录下执行如下脚本***
- ***如下脚本中的环境变量需根据具体情况修改***

```sh
# 如下环境变量根据具体情况修改

export ELASTICSEARCH_URL=http://127.0.0.1:9200 #elasticsearch的地址
export FILEBEAT_INDEX=filebeat-* # 索引名称前缀
export FILEBEAT_AUTH_USERNAME=xxxx #filebeat账号
export FILEBEAT_AUTH_PASSWORD=xxx #filebeat密码

# 创建 filebeat_writer角色
# 索引：filebeat-*
curl -X POST "${ELASTICSEARCH_URL}/_xpack/security/role/filebeat_writer?pretty" \
  -u ${AUTH_USERNAME}:${AUTH_PASSWORD} \
  -H 'Content-Type: application/json' \
  -d'
  {
    "cluster": ["manage_index_templates","monitor","manage_ingest_pipelines"], 
    "indices": [
      {
        "names": [ "'"${FILEBEAT_INDEX}"'" ], 
        "privileges": ["write","create_index"]
      }
    ]
  }'

# OR 启用ILM的时权限设置
# 索引：filebeat-*
curl -X PUT "${ELASTICSEARCH_URL}/_xpack/security/role/filebeat_writer?pretty" \
  -u ${AUTH_USERNAME}:${AUTH_PASSWORD} \
  -H 'Content-Type: application/json' \
  -d'
  {
    "cluster": ["manage_index_templates","monitor","manage_ingest_pipelines","manage_index_templates","manage_ilm"], 
    "indices": [
      {
        "names": [ "'"${FILEBEAT_INDEX}"'" ], 
        "privileges": ["write","create_index","manage"]
      }
    ]
  }'

# 创建 filebeat_ilm 角色
curl -X POST "${ELASTICSEARCH_URL}/_xpack/security/role/filebeat_ilm?pretty" \
  -u ${AUTH_USERNAME}:${AUTH_PASSWORD} \
  -H 'Content-Type: application/json' \
  -d'
  {
    "cluster": ["manage_ilm"],
    "indices": [
      {
        "names": [ "'"${FILEBEAT_INDEX}"'","shrink-filebeat-*"],
        "privileges": ["write","create_index","manage","manage_ilm"]
      }
    ]
  }'

# 创建 filebeat用户
curl -X POST "${ELASTICSEARCH_URL}/_xpack/security/user/${FILEBEAT_AUTH_USERNAME}?pretty" \
  -u ${AUTH_USERNAME}:${AUTH_PASSWORD} \
  -H 'Content-Type: application/json' \
  -d'
  {
    "password" : "'"${FILEBEAT_AUTH_PASSWORD}"'",
    "roles" : [ "filebeat_writer","kibana_user"],
    "full_name" : "'"${FILEBEAT_AUTH_USERNAME}"'"
  }'
```


### 创建ILM策略

```sh
export AUTH_USERNAME=xxx #用于链接
export AUTH_PASSWORD=xxx #用于链接

curl -X PUT "${ELASTICSEARCH_URL}/_ilm/policy/beats-default-policy" \
  -u ${AUTH_USERNAME}:${AUTH_PASSWORD} \
  -H 'Content-Type: application/json' \
  -d'
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_size": "800mb",
            "max_age": "3d"
          }
        }
      },
      "delete": {
        "min_age": "1h",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}'
```

#### 设置ILM策略周期

由于索引生命周期策略默认是10分钟检查一次符合策略的索引，因此在这10分钟内索引中的数据可能会超出指定的阈值。例如在步骤二：创建ILM策略时，设置max_docs为100，但doc数量在超过100后才触发索引滚动更新，此时可通过修改indices.lifecycle.poll_interval参数来控制检查频率，使索引在阈值范围内滚动更新。

```sh
curl -X PUT "${ELASTICSEARCH_URL}/_cluster/settings" \
  -u ${AUTH_USERNAME}:${AUTH_PASSWORD} \
  -H 'Content-Type: application/json' \
  -d'{
  "transient": {
    "indices.lifecycle.poll_interval":"1m"
  }
}'
```

## 安装

### 通过Ansible安装

- ***执行环境：类Unix系统；软件：ansible,ansible-galaxy***
- ***可以在任何目录下执行如下脚本***
- ***如下脚本中的环境变量需根据具体情况修改***

```sh
# 数据存放目录 待验证
# export FILEBEAT_CONFIG_DIR=/data/beats/filebeat/config
# export FILEBEAT_CONFIG_PATH=/data/beats/filebeat/config/filebeat.yml
# export FILEBEAT_DATA_DIR=/data/beats/filebeat/data
# export FILEBEAT_LOG_DIR=/data/beats/filebeat/log

# 如下环境变量根据具体情况修改

# 日志存放目录
export FILEBEAT_VERSION=6.8.4 # filebeat版本
export AUTH_USERNAME=xxx #用于链接kibana
export AUTH_PASSWORD=xxx #用于链接kibana
export FILEBEAT_AUTH_USERNAME=xxx #上一步骤创建的用户账号
export FILEBEAT_AUTH_PASSWORD=xxx #上一步骤创建的用户密码
export ELASTICSEARCH_URL=http://127.0.0.1:9200 #elasticsearch的地址
export FILEBEAT_HOST=127.0.0.1 # 安装filebeat的主机 可以添加多个
export KIBANA_URL=http://127.0.0.1:5601 # kibana地址
export FILEBEAT_PACKAGE_YRL="https://mirrors.huaweicloud.com/filebeat/${FILEBEAT_VERSION}/filebeat-${FILEBEAT_VERSION}-x86_64.rpm" #filebeatrpm包下载地址（使用华为的镜像）

export SSH_CONFIG="${PWD}/ssh-config" # ssh配置文件
export CLUSTER_NAME=filebat-test #集群名称 用户日志tag
export BASTION_SERVER=xxx.xxx.xxx.xxx #跳板机登陆服务器IP
export BASTION_NDENTIFY_FILE="${PWD}/bastion.pem" #跳板机登陆ssh密钥
export SERVICE_BUILT_IN_SERVER_HOST=127.0.0.1 #自建服务服务器地址
export SERVICE_BUILT_IN_SERVER_INDENTIFY_FILE="${PWD}/service-built-in.pem" #自建服务服务器登陆ssh密钥

# redis服务地址配置
export REDIS_SERVER_HOST=xxx.xxx.xxx.xxx #server地址
export REDIS_SERVER_PORT=7001 #端口

# 安装role
ansible-galaxy install elastic.beats,v7.9.2

cat << EOF > ${SSH_CONFIG}
Host ${CLUSTER_NAME}-serivce
     HostName ${SERVICE_BUILT_IN_SERVER_HOST}
     User root
     Port 22
     ProxyCommand ssh -W %h:%p ${CLUSTER_NAME}-bastion
     IdentityFile ${SERVICE_BUILT_IN_SERVER_INDENTIFY_FILE}

Host ${CLUSTER_NAME}-bastion
     HostName ${BASTION_SERVER}
     User root
     Port 22
     IdentityFile ${BASTION_NDENTIFY_FILE}
EOF

# ansible host使用
cat << EOF > filebeat
[filebeats]
${CLUSTER_NAME}-serivce ansible_ssh_common_args="-F ${SSH_CONFIG}"
EOF

cat << EOF > filebeat-playbook.yml
- name: filebeat
  hosts: filebeats
  roles:
    - elastic.beats
  vars:
    beats_version: ${FILEBEAT_VERSION}
    beat: filebeat #必须
    use_repository: false #国内yum安装受阻，所以使用rpm包形式安装
    custom_package_url: "${FILEBEAT_PACKAGE_YRL}"
    beat_conf:
      filebeat:
        config:
          modules:
            path: \${path.config}/modules.d/*.yml #必须
            reload.enabled: true #可选
        fields_under_root: true
        fields:
          cluster: "${CLUSTER_NAME}"
          sourcetype: "base-service"  
        inputs:
          - type: log #rabbitmq日志
            enabled: true # 必须启用
            paths: # 日志文件路径数组
              - /var/log/rabbitmq/*.log
            tags: ["common-service","rabbitmq","${CLUSTER_NAME}"]
            processors:
              - add_cloud_metadata: ~ # 云提供商元数据
              - add_host_metadata: ~ # 主机元数据
          - type: log #redis日志
            paths: # redis日志文件路径数组
              - /var/log/redis.log
            tags: ["common-service","redis","${CLUSTER_NAME}"] 
            processors:
              - add_cloud_metadata: ~ # 云提供商元数据
              - add_host_metadata: ~ # 主机元数据
          - type: redis #redis慢查询日志
            hosts: ["${REDIS_SERVER_HOST}:${REDIS_SERVER_PORT}"]
            #password: "${redis_pwd}"
            tags: ["common-service","redis","slowlogs","${CLUSTER_NAME}"] 
            processors:
              - add_cloud_metadata: ~ # 云提供商元数据
              - add_host_metadata: ~ # 主机元数据
        #modules:
        #  - module: redis
      setup:
        kibana:
          host: "${KIBANA_URL}"
          username: ${AUTH_USERNAME}
          password: ${AUTH_PASSWORD}
        dashboards:
          enable: true #启动kibana filebeat dashboard
    output_conf:
      elasticsearch:
        hosts: ["${ELASTICSEARCH_URL}"]
        username: ${FILEBEAT_AUTH_USERNAME}
        password: ${FILEBEAT_AUTH_PASSWORD}
        ilm.enabled: true
        ilm.rollover_alias: "filebeat"
        ilm.pattern: "{now/d}-000001"
EOF

# 运行playbook安装filebeat
ansible-playbook -i filebeat filebeat-playbook.yml
```

#### 输出到`logstash`

将 `output_conf`修改如下

- 127.0.0.1：logstash地址
- 5044：logstash beat端口

```yaml
output_conf:
  logstash:
    hosts: ["127.0.0.1:5044"]
```

### 在kubernetes环境中安装

#### 参考链接

- https://www.elastic.co/guide/en/beats/filebeat/6.8/running-on-kubernetes.html

#### 安装

- ***执行环境：类Unix系统；软件：sed kubectl***
- ***可以在任何目录下执行如下脚本***
- ***如下脚本中的环境变量需根据具体情况修改***

```sh
# 设置模版中的参数值
export ELASTICSEARCH_HOST=127.0.0.1
export ELASTICSEARCH_PORT=9200
export ELASTICSEARCH_USERNAME=xxx
export ELASTICSEARCH_PASSWORD=xxx
export CLUSTER_NAME=filebat-test #集群名称 用户日志tag

# 替换模版中占位符为真正值
sed 's/%CLUSTER_NAME%/'"${CLUSTER_NAME}"'/g;s/%ELASTICSEARCH_HOST%/'"${ELASTICSEARCH_HOST}"'/g;s/%ELASTICSEARCH_PORT%/'"${ELASTICSEARCH_PORT}"'/g;s/%ELASTICSEARCH_USERNAME%/'"${ELASTICSEARCH_USERNAME}"'/g;s/%ELASTICSEARCH_PASSWORD%/'"${ELASTICSEARCH_PASSWORD}"'/g' filebeat-kubernetes-v1.19.7.yml.template > filebeat-kubernetes-v1.19.7.yml

# 部署
kubectl create -f filebeat-kubernetes-v1.19.7.yml
```

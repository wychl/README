# 部署

## 安装

- ***执行环境：类Unix系统；软件：ansible,ansible-galaxy***
- ***可以在任何目录下执行如下脚本***
- ***如下脚本中的环境变量需根据具体情况修改***

```sh
# 数据存放目录
export DATA_DIR=/data/elasticsearch/data
# 日志存放目录
export LOG_DIR=/data/elasticsearch/log
export VERSION=6.8.4
export HTTP_HOST=0.0.0.0
export HTTP_PORT=9200
export TRANSPORT_HOST=0.0.0.0
export TRANSPORT_PORT=9300
export AUTH_USERNAME=admin
export AUTH_PASSWORD=admin

# 安装role
# ansible-galaxy install elastic.elasticsearch,v7.12.0
ansible-galaxy install elastic.elasticsearch,v7.12.0

# 创建host文件
cat << EOF >> elasticsearch-host
[elastic]
127.0.0.1 ansible_ssh_user=root
EOF

# 创建playbook
cat << EOF > elasticsearch.yml
- name: elasticsearch
  hosts: elastic
  roles:
    - role: elastic.elasticsearch
  vars:
    # 堆内存容量
    es_heap_size: "1g"
    es_version: ${VERSION}
    es_data_dirs:
      - ${DATA_DIR}
    es_log_dir: ${LOG_DIR}
    es_max_open_files: 65536 # 文件描述符数量
    # 使用rpm包安装
    es_use_repository: false
    es_custom_package_url: "https://mirrors.huaweicloud.com/elasticsearch/${VERSION}/elasticsearch-${VERSION}.rpm"
    es_api_basic_auth_username: ${AUTH_USERNAME}
    es_api_basic_auth_password: ${AUTH_PASSWORD}
    es_config:
      cluster.name: "log-cluster" #集群名称
      # 服务发现和集群形成设置
      #cluster.initial_master_nodes: "log-cluster02" # 初始主节点
      #discovery.seed_hosts: "log-cluster:9300" # 服务发现种子主机
      node.name: log-node-01
      http.host: ${HTTP_HOST}
      http.port: ${HTTP_PORT}
      node.data: true
      node.master: true
      transport.host: ${TRANSPORT_HOST}
      transport.port: ${TRANSPORT_PORT}
      bootstrap.memory_lock: false
    es_plugins:
     - plugin: ingest-attachment
    es_users:
      file:
        xxx:
          password: xxx
          roles:
            - admin
        testUser:
          password: testUser!
          roles:
            - power_user
            - user
    es_roles:
      file:
        admin:
          cluster:
            - all
          indices:
            - names: '*'
              privileges:
                - all
        power_user:
          cluster:
            - monitor
          indices:
            - names: '*'
              privileges:
                - all
        user:
          indices:
            - names: '*'
              privileges:
                - read
EOF

# 运行playbook安装elasticsearch
ansible-playbook -i elasticsearch-host elasticsearch.yml
```

### 测试

```sh
# 配置环境变量
export ELASTICSEARCH_HOST=127.0.0.1
```

[测试脚本](../script/elasticsearch.sh)
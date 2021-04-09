
# kibana

## 安装

### docker run

```sh
export KIBANE_DATA=$PWD/kibana/data
export KIBANE_CONFIG=$PWD/kibana/config
docker run -d \
  --name kibana \
  -p 5601:5601  \
  -d \
  --restart=always \
  -v ${KIBANE_CONFIG}:/usr/share/kibana/config \
  -v ${KIBANE_DATA}:/usr/share/kibana/data \
  docker.elastic.co/kibana/kibana:6.8.4
```

### binary run

- ***执行环境：类Unix系统；软件：ansible,ansible-galaxy***
- ***在当前目录下执行如下脚本***
- ***如下脚本中的环境变量需根据具体情况修改***

```sh
# 安装role
export KIBANA_CONFIG_DIR=/data/kibana/config
export KIBANA_CONFIG=/data/kibana/config/kibana.yml
export KIBANA_VERSION=6.8.4
export KIBANA_SERVER_HOST=0.0.0.0
export KIBANA_SERVER_PORT=5601
export KIBANA_ELASTICSEARCH_URL=http://127.0.0.1:9200
export ES_USERNAME=xxx
export ES_PASSWORD=xxx
export KIBANA_DATA_DIR=/data/kibana/data

cat << EOF >> elastic
[kibana]
127.0.0.1 ansible_ssh_user=root
EOF

# 创建playbook
cat << EOF > kibana.yml
- name: kibana
  hosts: kibana
  roles:
    - kibana
  vars:
    kibana_version: ${KIBANA_VERSION}
    kibana_server_port: ${KIBANA_SERVER_PORT}
    kibana_server_host: ${KIBANA_SERVER_HOST}
    kibana_elasticsearch_username: ${ES_USERNAME}
    kibana_elasticsearch_password: ${ES_PASSWORD}
    kibana_config_dir: ${KIBANA_CONFIG_DIR}
    kibana_config_file_path: ${KIBANA_CONFIG}
    kibana_elasticsearch_url: ${KIBANA_ELASTICSEARCH_URL}
    # 使用rpm包安装
    kibana_use_repository: false
    kibana_custom_package_url: "https://mirrors.huaweicloud.com/kibana/${VERSION}/kibana-${VERSION}-x86_64.rpm"
    kibana_data_dir: ${KIBANA_DATA_DIR}
EOF

# 运行playbook安装elasticsearch
ansible-playbook -i elastic kibana.yml
```

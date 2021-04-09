# logstash

## 参考链接

- https://github.com/geerlingguy/ansible-role-logstash
- https://www.elastic.co/guide/en/beats/filebeat/6.8/configuring-ssl-logstash.html
- https://www.elastic.co/guide/en/elasticsearch/reference/6.8/certutil.html

## 安装

- ***执行环境：类Unix系统；软件：ansible,ansible-galaxy***
- ***在当前目录下执行如下脚本***
- ***如下脚本中的环境变量需根据具体情况修改***

```sh
# 如下环境变量根据具体情况修改
export LOGSTASH_MAIN_VERSION=6.x # logstash版本
export LOGSTASH_VERSION=6.8.4 # logstash版本
# export LOGSTASH_DIR=/data/logstash/share
export AUTH_USERNAME=xxx #用于链接elasticsearch账号
export AUTH_PASSWORD=xxx #用于链接elasticsearch密码
export LOGSTASH_AUTH_USERNAME=logstash #上一步骤创建的用户账号
export LOGSTASH_AUTH_PASSWORD=logstash #上一步骤创建的用户密码
export ELASTICSEARCH_URL=http://127.0.0.1:9200 #elasticsearch的地址
export LOGSTASH_HOST=127.0.0.1 # 安装logstash的主机 可以添加多个
export LOGSTASH_PACKAGE_YRL="https://mirrors.huaweicloud.com/logstash/${LOGSTASH_VERSION}/logstash-${LOGSTASH_VERSION}.rpm" #logstash rpm包下载地址（使用华为的镜像）
export LOGSTASH_LISTEN_PORT_BEATS=5044 # logstash监听端口，beats将数据发送这个端口，然后由logstash转发或者处理
export LOGSTASH_GEM_SOUECE="https://ruby.taobao.org/" #设置gem下载源；解决插件无法下载问题

cat << EOF > logstash
[logstash]
${LOGSTASH_HOST} ansible_ssh_user=root
EOF

cat << EOF > logstash-playbook.yml
- name: logstash
  hosts: logstash
  roles:
    - logstash
  vars:
    logstash_version: ${LOGSTASH_MAIN_VERSION}
    logstash_elasticsearch_hosts: ["${ELASTICSEARCH_URL}"]
    logstash_use_repository: false
    logstash_custom_package_url: "${LOGSTASH_PACKAGE_YRL}"
    logstash_elasticsearch_username: ${AUTH_USERNAME}
    logstash_elasticsearch_password: ${AUTH_PASSWORD}
    logstash_listen_port_beats: ${LOGSTASH_LISTEN_PORT_BEATS}
    logstash_gem_source: ${LOGSTASH_GEM_SOUECE}
    logstash_install_plugins: # 安装的插件
      - logstash-input-beats
EOF

# 运行playbook安装filebeat
ansible-playbook -i logstash logstash-playbook.yml
```

## debug

```sh
/usr/share/logstash/bin/logstash --path.settings /etc/logstash -e

# 查看日志
journalctl -fu logstash.service
```
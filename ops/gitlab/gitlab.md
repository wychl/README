# gitlab

## 参看链接

- https://docs.gitlab.com/omnibus/docker/

## gitlab服务器环境描述

- gitlab-ce 版本:`12.4.3`
- docker部署 版本:`19.03.15-3.el8`
- 阿里云对象存储：gitlab备份文件存储
- 硬件
  - 数据盘：70GB
  - 系统盘：40GB
  - 内存：16GB
  - CPU：2核
- https访问
- ntp

## 服务器环境配置

### 挂载阿里云磁盘

***根据需要执行***

```sh
#查看分区
df -h
# 输出内容如下
# 文件系统        容量  已用  可用 已用% 挂载点
# devtmpfs        7.7G     0  7.7G    0% /dev
# tmpfs           7.7G     0  7.7G    0% /dev/shm
# tmpfs           7.7G  472K  7.7G    1% /run
# tmpfs           7.7G     0  7.7G    0% /sys/fs/cgroup
# /dev/vda1        40G  2.1G   36G    6% /
# tmpfs           1.6G     0  1.6G    0% /run/user/0


# 查看磁盘列表
fdisk -l
# 输出内容如下
# 磁盘 /dev/vda：42.9 GB, 42949672960 字节，83886080 个扇区
# Units = 扇区 of 1 * 512 = 512 bytes
# 扇区大小(逻辑/物理)：512 字节 / 512 字节
# I/O 大小(最小/最佳)：512 字节 / 512 字节
# 磁盘标签类型：dos
# 磁盘标识符：0x000bb9c1

# 设备 Boot      Start         End      Blocks   Id  System
# /dev/vda1   *        2048    83886046    41941999+  83  Linux

# 磁盘 /dev/vdb：75.2 GB, 75161927680 字节，146800640 个扇区
# Units = 扇区 of 1 * 512 = 512 bytes
# 扇区大小(逻辑/物理)：512 字节 / 512 字节
# I/O 大小(最小/最佳)：512 字节 / 512 字节


# 分区
# 输入n，p，1，回车，回车，wq，保存退出。
fdisk /dev/vdb
# 输出内容如下
# 欢迎使用 fdisk (util-linux 2.23.2)。

# 更改将停留在内存中，直到您决定将更改写入磁盘。
# 使用写入命令前请三思。

# Device does not contain a recognized partition table
# 使用磁盘标识符 0xfc08face 创建新的 DOS 磁盘标签。

# 命令(输入 m 获取帮助)：n
# Partition type:
#    p   primary (0 primary, 0 extended, 4 free)
#    e   extended
# Select (default p): p
# 分区号 (1-4，默认 1)：1
# 起始 扇区 (2048-146800639，默认为 2048)：
# 将使用默认值 2048
# Last 扇区, +扇区 or +size{K,M,G} (2048-146800639，默认为 146800639)：
# 将使用默认值 146800639
# 分区 1 已设置为 Linux 类型，大小设为 70 GiB

# 命令(输入 m 获取帮助)：wq
# The partition table has been altered!

# Calling ioctl() to re-read partition table.
# 正在同步磁盘。


# 格式化分区
mkfs.ext4 /dev/vdb1
# 创建目录挂载
echo '/dev/vdb1  /data ext4    defaults    0  0' >> /etc/fstab
# 执行mount挂载操作
mkdir -p /data
mount /dev/vdb1 /data


#查看分区
df -h
# 文件系统        容量  已用  可用 已用% 挂载点
# devtmpfs        7.7G     0  7.7G    0% /dev
# tmpfs           7.7G     0  7.7G    0% /dev/shm
# tmpfs           7.7G  476K  7.7G    1% /run
# tmpfs           7.7G     0  7.7G    0% /sys/fs/cgroup
# /dev/vda1        40G  2.1G   36G    6% /
# tmpfs           1.6G     0  1.6G    0% /run/user/0
# /dev/vdb1        69G   52M   66G    1% /data
```

### 系统

```sh
# 时区
timedatectl set-timezone Asia/Shanghai

# 设置hostbname
hostnamectl set-hostname gitlab

# 更新系统软件
yum update
```

### docker安装和配置

#### 卸载老版本

***根据需要执行***

```sh
sudo yum remove docker \
    docker-client \
    docker-client-latest \
    docker-common \
    docker-latest \
    docker-latest-logrotate \
    docker-logrotate \
    docker-engine
```

#### 设置yum

```sh
sudo yum install -y yum-utils
sudo yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

#### 安装

```sh
# 查看可安装的docker版本
yum list docker-ce --showduplicates | sort -r

# 执行安装
sudo yum install \
  docker-ce \
  docker-ce-cli \
  containerd.io

# 也可以执行版本安装，这里选择的是安装最新版本
# sudo yum install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io
```

#### 创建docker用户组

```sh
sudo groupadd docker
sudo usermod -aG docker $USER
```

#### 启动

```sh
# 启动docker服务
sudo systemctl start docker
# 查看启动状态
sudo systemctl status docker
# 配置自启动
sudo systemctl enable docker
sudo systemctl enable containerd
# 测试
docker run --rm hello-world
```

#### 配置数据存储路径

编辑文件 `/etc/docker/daemon.json`的`data-root`字段（值为：docker数据存储目录）

- 操作示例如下

```sh
# 数据存储目录为/data/docker
mkdir -p  /data/docker
vi /etc/docker/daemon.json
```

```json
{
 "data-root": "/data/docker"
}
```

- 重启docker服务

```sh
sudo systemctl stop docker
sudo systemctl start docker
sudo systemctl status docker
```

## gitlab搭建

### 配置环境变量

```sh
export GITLAB_HOME=/data/gitlab #gitlab数据文件目录（配置、数据和日志）
export GITLAB_CONTAINER_NAME=gitlab # 容器名称
export HOSTNAME=xxx.com # gitlab域名
mkdir -p ${GITLAB_HOME} # 创建数据文件目录
```

### 启动服务

```sh
docker run --detach \
  --hostname ${HOSTNAME} \
  --publish 8000:8000 --publish 2202:22 \
  --name ${GITLAB_CONTAINER_NAME} \
  --restart always \
  --volume ${GITLAB_HOME}/secret:/secret/gitlab/backups \
  --volume ${GITLAB_HOME}/config:/etc/gitlab \
  --volume ${GITLAB_HOME}/logs:/var/log/gitlab \
  --volume ${GITLAB_HOME}/data:/var/opt/gitlab \
  gitlab/gitlab-ce:12.4.3-ce.0
```

### 修改默认配置

***配置文件路径 `/data/gitlab/gitlab.rb`***

```ruby
#路径：/etc/gitlab/gitlab.rb

#######################################配置http端口##############################################
external_url "http://xxx.com:8000" 
#######################################配置http端口##############################################



#######################################配置ssh协议端口##############################################
gitlab_rails['gitlab_shell_ssh_port']=2202
#######################################配置ssh协议端口##############################################


#######################################oss配置##############################################
# 备份文件oss存储
# 参考文件：https://github.com/fog/fog-aliyun
gitlab_rails['backup_upload_connection'] = {
'provider' => 'aliyun',
'aliyun_accesskey_id' => '有权限访问存储桶的用户key',
'aliyun_accesskey_secret' => '有权限访问存储桶的密钥',
'aliyun_oss_endpoint' => 'http://oss-cn-shanghai-internal.aliyuncs.com',
'aliyun_oss_bucket' => 'bucket',    #OSS名称
'aliyun_oss_location' => 'shanghai'      #oss地域
}
gitlab_rails['backup_upload_remote_directory'] = 'gitlab'    #存储gitlab备份的桶子目录
gitlab_rails['backup_multipart_chunk_size'] = 104857600 # 分片大小
#######################################oss配置##############################################


#######################################邮箱配置##############################################
# SMTP 方式发送邮件
gitlab_rails['smtp_enable'] = true
gitlab_rails['smtp_address'] = "smtp.exmail.qq.com"
gitlab_rails['smtp_port'] = 465

# 邮箱账号和密码
gitlab_rails['smtp_user_name'] = "xxxx@gemii.cc"
gitlab_rails['smtp_password'] = "xxxx"
gitlab_rails['smtp_domain'] = "gemii.cc"
gitlab_rails['smtp_authentication'] = "login"

# If your SMTP server does not like the default 'From: gitlab@localhost' you
# can change the 'From' with this setting.
# gitlab 发件人
gitlab_rails['gitlab_email_enabled'] = true
gitlab_rails['gitlab_email_from'] = 'gitlab@example.com'
gitlab_rails['gitlab_email_reply_to'] = 'noreply@example.com'
#######################################邮箱配置##############################################

#######################################时区配置##############################################
gitlab_rails['time_zone'] = 'Asia/Shanghai'
#######################################时区配置##############################################

```

### 应用新的gitlab配置

```sh
# 将配置文件放到服务器
scp gitlab.rb user@host:/data/gitlab/config

# 重启gitlab
docker stop ${GITLAB_CONTAINER_NAME}
docker start ${GITLAB_CONTAINER_NAME}
```

## gitlab安装之后的

### 注册限制

- Sign-up enabled 是否开启注册
- Require admin approval for new sign-ups 是否需要审批
- Send confirmation email on sign-up 是否发送确认邮件
- Minimum password length 密码最小长度
- Email restrictions for sign-ups 注册时的邮箱地址限制

## Gitlab迁移

### 迁移参考文档

- https://docs.gitlab.com/omnibus/settings/backups.html#backup-and-restore-omnibus-gitlab-configuration
- https://docs.gitlab.com/ce/raketasks/backup_restore.html#restore-for-omnibus-installations
- https://www.codenong.com/js40a1e26b2e1d/

### 备份与恢复操作步骤-示例

#### 数据备份

```sh
docker exec -t ${GITLAB_CONTAINER_NAME} gitlab-backup
```

### 恢复

```sh
#数据备份的文件放到备份目录
mv 1614703726_2021_03_02_12.4.3_gitlab_backup.tar ${GITLAB_HOME}/data/backups/
# 停止链接数据库的进程
docker exec -it ${GITLAB_CONTAINER_NAME} gitlab-ctl stop unicorn
docker exec -it ${GITLAB_CONTAINER_NAME} gitlab-ctl stop puma
docker exec -it ${GITLAB_CONTAINER_NAME} gitlab-ctl stop sidekiq
# Verify that the processes are all down before continuing
docker exec -it ${GITLAB_CONTAINER_NAME} gitlab-ctl status
# 更改容器内部文件权限
docker exec -it ${GITLAB_CONTAINER_NAME} chown git:git /var/opt/gitlab/backups
# 运行恢复命令 加入文件名为 1614703726_2021_03_02_12.4.3_gitlab_backup.tar 则恢复命令如下
docker exec -it ${GITLAB_CONTAINER_NAME} gitlab-backup restore BACKUP=1615288278_2021_03_09_12.4.3
# 重启gitlab服务
docker restart ${GITLAB_CONTAINER_NAME}
# 查看启动日志
docker logs -f ${GITLAB_CONTAINER_NAME}
# 检测GitLab是否正常运行
docker exec -it ${GITLAB_CONTAINER_NAME} gitlab-rake gitlab:check SANITIZE=true
```

### 使用`conntab`定时备份gitlab数据

```sh
# 创建备份脚本
export BACKUP_SHELL=$PWD/backup.sh
cat << EOF > ${BACKUP_SHELL}
#!/bin/bash
docker exec gitlab gitlab-backup
EOF

# 给予脚本执行权限

chmod +x ${BACKUP_SHELL}

# 创建定时任务
crontab -e

# 输入一下内容
# 0 2 * * * /path/to/backup-script
```

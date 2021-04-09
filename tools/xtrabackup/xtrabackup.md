# 数据库备份和恢复

## 参考文档

- https://www.percona.com/doc/percona-xtrabackup/2.4/index.html
- https://www.cnblogs.com/linuxk/p/9372990.html
- https://help.aliyun.com/knowledge_detail/41817.html?spm=5176.11065259.1996646101.searchclickresult.53d420cclqekK3

## 描述

- MySQL 5.6及之前的版本需要安装 Percona XtraBackup 2.3，安装指导请参见官方文档Percona XtraBackup 2.3。
- MySQL 5.7版本需要安装 Percona XtraBackup 2.4，安装指导请参见官方文档Percona XtraBackup 2.4。
- MySQL 8.0版本需要安装 Percona XtraBackup 8.0，安装指导请参见官方文档Percona XtraBackup 8.0。

## 工具介绍

### 原理

![备份原理图](./images/backup.png)

1. innobackupex启动后，会先fork一个进程，用于启动xtrabackup，然后等待xtrabackup备份ibd数据文件；
2. xtrabackup在备份innoDB数据是，有2种线程：redo拷贝线程和ibd数据拷贝线程。xtrabackup进程开始执行后，会启动一个redo拷贝的线程，用于从最新的checkpoint点开始顺序拷贝redo.log；再启动ibd数据拷贝线程，进行拷贝ibd数据。这里是先启动redo拷贝线程的。在此阶段，innobackupex进行处于等待状态（等待文件被创建）
3. xtrabackup拷贝完成ibd数据文件后，会通知innobackupex（通过创建文件），同时xtrabackup进入等待状态（redo线程依旧在拷贝redo.log）
4. innobackupex收到xtrabackup通知后哦，执行FLUSH TABLES WITH READ LOCK（FTWRL），取得一致性位点，然后开始备份非InnoDB文件（如frm、MYD、MYI、CSV、opt、par等格式的文件），在拷贝非InnoDB文件的过程当中，数据库处于全局只读状态。
5. 当innobackup拷贝完所有的非InnoDB文件后，会通知xtrabackup，通知完成后，进入等待状态；
6. xtrabackup收到innobackupex备份完成的通知后，会停止redo拷贝线程，然后通知innobackupex，redo.log文件拷贝完成；
7. innobackupex收到redo.log备份完成后，就进行解锁操作，执行：UNLOCK TABLES；
8. 最后innbackupex和xtrabackup进程各自释放资源，写备份元数据信息等，innobackupex等xtrabackup子进程结束后退出。

### xtrabackup

是用于热备innodb，xtradb表中数据的工具，不能备份其他类型的表，也不能备份数据表结构；

### innobackupex

是将xtrabackup进行封装的perl脚本，提供了备份myisam表的能力。

### 常用选项

- host 指定主机
- user 指定用户名
- password 指定密码
- port 指定端口
- databases 指定数据库
- incremental 创建增量备份
- incremental-basedir 指定包含完全备份的目录
- incremental-dir   指定包含增量备份的目录
- apply-log 对备份进行预处理操作，一般情况下，在备份完成后，数据尚且不能用于恢复操作，因为备份的数据中可能会包含尚未提交的事务或已经提交但尚未同步至数据文件中的事务。因此，此时数据文件仍处理不一致状态。“准备”的主要作用正是通过回滚未提交的事务及同步已经提交的事务至数据文件也使得数据文件处于一致性状态。
- redo-only 不回滚未提交事务
- copy-back 恢复备份目录

## 操作环境

- 操作系统 RedHat
- MySQL 5.7
  - 版本 5.7
  - 备份用户和权限 'xtrabackup'@'localhost'
- XtraBackup 2.4

### MySQL

#### 启动MySQL

```sh
export MYSQL_DATA=/data/mysql
mkdir -p ${MYSQL_DATA}
docker run --name mysql \
    -e MYSQL_ROOT_PASSWORD=root \
    -p 3306:3306 \
    -v ${MYSQL_DATA}:/var/lib/mysql \
    -d \
    mysql:5.7 \
    --character-set-server=utf8mb4 \
    --collation-server=utf8mb4_unicode_ci
```

#### 创建备份用户和配置权限

- https://www.percona.com/doc/percona-xtrabackup/2.4/using_xtrabackup/privileges.html
- https://dev.mysql.com/doc/refman/5.7/en/create-user.html
- https://dev.mysql.com/doc/refman/5.7/en/account-names.html
- https://dev.mysql.com/doc/refman/5.7/en/privileges-provided.html

```sql
CREATE USER 'xtrabackup'@'localhost' IDENTIFIED BY 'xtrabackup'; #创建用户
GRANT RELOAD, LOCK TABLES, PROCESS, REPLICATION CLIENT ON *.* TO 'xtrabackup'@'localhost'; #刷新、锁定表、用户查看服务器状态
FLUSH PRIVILEGES; #刷新授权表
REVOKE ALL PRIVILEGES,GRANT OPTION FROM 'xtrabackup';　　# 回收此用户所有权限
```

## 安装

```sh
wget https://www.percona.com/downloads/XtraBackup/Percona-XtraBackup-2.4.4/binary/redhat/7/x86_64/percona-xtrabackup-24-2.4.4-1.el7.x86_64.rpm
yum localinstall percona-xtrabackup-24-2.4.4-1.el7.x86_64.rpm
```

## 使用

### 备份

#### 完全备份

##### 创建数据备份目录

```sh
export BACKUP_DIR=/backup/mysql
mkdir -p ${BACKUP_DIR}
cd ${BACKUP_DIR}
```

##### 创建备份

```sh
# 不指定数据库
xtrabackup --uxtrabackup -pxtrabackup --host=127.0.0.1 --backup --datadir=${MYSQL_DATA} --target-dir=${BACKUP_DIR} 

# 执行数据库
xtrabackup --uxtrabackup -pxtrabackup --host=127.0.0.1 --backup --datadir=${MYSQL_DATA} --database=test --target-dir=${BACKUP_DIR}  


innobackupex --uxtrabackup -pxtrabackup /path/to/backup/dir/
```

##### 备份

```sh
xtrabackup --prepare --target-dir=${BACKUP_DIR}
```

### 恢复

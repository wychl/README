# docker运行mariadb

## 创建自定义配置文件和数据库数据存放路径

- 创建配置文件路径 `sudo mkdir /mariadb/config -p`
- 创建数据文件存放路径 `sudo mkdir /mariadb/data -p`

## 将utf8mb4.cnf拷贝到`/mariadb/config`目录下

## 启动容器

```shell
sudo docker run --name mariadb -v /mariadb/config:/etc/mysql/conf.d -v /mariadb/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=0000 -p 3306:3306 -d mariadb
```

## 查看编码方式

```shell
show variables like 'character%';
```
## 将shadowsocks转化为http代理

### 安装polipo和设置代理:
```shell
sudo apt-get install polipo
sudo service polipo stop
sudo polipo socksParentProxy=localhost:1080 &
sudo service polipo start
```
### 使用代理:

```shell
http_proxy=http://localhost:8123 apt-get update

http_proxy=http://localhost:8123 curl www.google.com

http_proxy=http://localhost:8123 wget www.google.com

git config --global http.proxy 127.0.0.1:8123
git clone https://github.com/xxx/xxx.git
git xxx
git xxx
git config --global --unset-all http.proxy
```

### 将shdowsocks协议转化为http协议(docker方式)
- 使用的是go版的shadowsocks客户端
- [shdowsocks配置文件](./docker/config.json.sample) 使用时将文件名改为config.json,并改为自己的shadowsocks的配置
- [Dockerfile](./docker/Dockerfile)
- 构建镜像 `docker build . -t socks`
- 运行容器 `docker run --name socks -d -p 8123:8123 -p 1080:1080 socks:latest`
- http协议端口为8123

# dnsmasq

## 参考文档

- https://hub.docker.com/r/jpillora/dnsmasq
- http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html

Linux 处理 DNS 请求时有个限制，在 resolv.conf 中最多只能配置三个域名服务器（nameserver）。作为一种变通方法,可以在 resolv.conf 文件中只保留 localhost 作为域名服务器，然后为外部域名服务器另外创建 resolv-file 文件。首先，为 dnsmasq 新建一个域名解析文件：

```sh
vi /etc/resolv.dnsmasq.conf
# Google's nameservers, for example
nameserver 8.8.8.8
nameserver 8.8.4.4
```

## run

```sh
docker run \
    --name dnsmasq \
    -d \
    -p 53:53/udp \
    -p 5380:8080 \
    -v ${PWD}/dnsmasq.conf:/etc/dnsmasq.conf \
     -v ${PWD}/resolv.dnsmasq.conf:/etc/resolv.dnsmasq.conf \
    --log-opt "max-size=100m" \
    -e "HTTP_USER=foo" \
    -e "HTTP_PASS=bar" \
    --restart always \
    jpillora/dnsmasq
```

## 测试

```sh
$ host myhost.company <docker-host>
Using domain server:
Name: <docker-host>
Address: <docker-host>#53
Aliases:

myhost.company has address 10.0.0.2
```

在配置好 dnsmasq 后，你需要编辑/etc/resolv.conf 让 DHCP 客户端首先将本地地址(localhost)加入 DNS 文件(/etc/resolv.conf)，然后再通过其他 DNS 服务器解析地址。配置好 DHCP 客户端后需要重新启动网络来使设置生效。

resolv.conf
一种选择是一个纯粹的 resolv.conf 配置。要做到这一点，才使第一个域名服务器在/etc/resolv.conf 中指向 localhost：

/etc/resolv.conf
nameserver 127.0.0.1

# External nameservers

...
现在，DNS 查询将首先解析 dnsmasq，只检查外部的服务器如果 DNSMasq 无法解析查询. dhcpcd, 不幸的是，往往默认覆盖 /etc/resolv.conf, 所以如果你使用 DHCP，这里有一个好主意来保护 /etc/resolv.conf,要做到这一点，追加 nohook resolv.conf 到 dhcpcd 的配置文件：

/etc/dhcpcd.conf
...
nohook resolv.conf
也可以保护您的 resolv.conf 不被修改：

# chattr +i /etc/resolv.conf

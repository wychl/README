# 技巧

## vim 访问远程文件

```shell
 vim scp://example/example.text
```

## 远程服务当本地用

通过 LocalForward 将本地端口上的数据流量通过 ssh 转发到远程主机的指定端口。感觉你是使用的本地服务，其实你使用的远程服务。如远程服务器上运行着 Postgres，端口 5432（未暴露端口给外部）。那么，你可以：
Host db
    HostName db.example.com
    LocalForward 5433 localhost:5432
当你连接远程主机时，它会在本地打开一个 5433 端口，并将该端口的流量通过 ssh 转发到远程服务器上的 5432 端口。
首先，建立连接：
$ ssh db
之后，就可以通过 Postgres 客户端连接本地 5433 端口：
$ psql -h localhost -p 5433 orders

## 绑定本地端口

N参数，表示只连接远程主机，不打开远程shell；T参数，表示不为这个连接分配TTY。这个两个参数可以放在一起用，代表这个SSH连接只用来传数据，不执行远程操作。

```shell
ssh -NT -D 8080  user@host
```
f参数，表示SSH连接成功后，转入后台运行。这样一来，你就可以在不中断SSH连接的情况下，在本地shell中执行其他操作。

```shell
ssh -f -D 8080  user@host
```

SSH会建立一个socket，去监听本地的8080端口。一旦有数据传向那个端口，就自动把它转移到SSH连接上面，发往远程主机。可以想象，

如果8080端口原来是一个不加密端口，现在将变成一个加密端口。

## 远程操作

SSH不仅可以用于远程主机登录，还可以直接在远程主机上执行操作

```shell
ssh  user@host ls
```

## 本地端口转发

## 远程端口转发

## 加密 Socks 通道

对于连接到各种不安全的无线网络上的笔记本电脑用户来说,这个是特别有用的！唯一所需要的就是一个一定程度上处于安全的地点的 SSH 服

务器，比如在家里或办公室。用动态的 DNS 服务 DynDNS 也可能是很有用的，这样你就不必记住你的 IP 了。

1. 开始连接
你只要执行这一个命令就能开始你的连接：

`$ ssh -TND 4711 user@host`

这里的 user 是你在 host 这台 SSH 服务器上的用户名。它会让你输入密码，然后你就能连上了。 N 表示不采用交互提示，而 D 表示

指定监听的本地端口（你可以使用任何你喜欢的数字），T 表示禁用伪 tty 分配。

加了 -v (verbose) 标志以后的输出可以让你能够验证到底连了哪个端口。

2. 配置你的浏览器(或其它程序)

如果你没有配置你的浏览器（或其他程序）使用这个新创建的 socks 隧道，上述步骤是无效的。由于当前版本的 SSH 支持 SOCKS4 和 SOCKS5，因此您可以使用其中任何一种。

对于 Firefox: Edit > Preferences > Advanced > Network > Connection > Setting: 

选中 Manual proxy configuration 单选框, 然后在 SOCKS host 里输入 localhost， 然后在后面那个框中输入你的端口号（本例中为 4711）。

Firefox 不会自动通过 socks 隧道发送 DNS 请求，这一潜在的隐私问题可以通过以下步骤来解决：

在 Firefox 地址栏中输入：about:config 。

搜索：network.proxy.socks_remote_dns

将该值设为 true。

重启浏览器。

对于 Chromium: 你可以将 SOCKS 设置设置为环境变量或命令行选项。我建议将下列函数之一加入到你的 .bashrc：

```text
        function secure_chromium {
            port=4711
            export SOCKS_SERVER=localhost:$port
            export SOCKS_VERSION=5
            chromium &
            exit
        }
        或者

        function secure_chromium {
            port=4711
            chromium --proxy-server="socks://localhost:$port" &
            exit
        }
```

现在打开终端然后输入：

`$ secure_chromium`

享受你的安全隧道吧！

## 通过跳板记链接服务器

- 场景

```shell
    ssh       ssh
A ------> B ------> C
    ^          ^
 using A's   using B's
 ssh key     ssh key

```

- 条件:

A is running ssh-agent;
A can access B;
B can access C;
A's ssh public key is present in B:~/.ssh/authorized_keys
B's ssh public key is present in C:~/.ssh/authorized_keys

- 在A端~/.ssh/config加如下配置

``` shell
Host C
    ProxyCommand ssh -o 'ForwardAgent yes' B 'ssh-add && nc %h %p'
```

- 测试

```sh
ssh C
```

## 端口转发

- 示例：
    通过本地mysql客户端链接，远端没有包露出来mysql服务器

- shell 命令

```sh
ssh user@host -L 3307:127.0.0.1:3306 -N
```

- 测试

```sh
mysql -h 127.0.0.1 -u user -p password -P 3307 
```

# ssh_config配置

## 位置

- `~/.ssh/config`  用户配置文件
- `/etc/ssh/ssh_config` 系统配置文件

## 示例

```sh
Host example                        #主机别名
    HostName example.com            # 指定主机地址
    User root                       # 指定用户名
    IdentityFile ~/.ssh/id_ecdsa    # 指定认证文件位置
    Port 22                         # 指定端口
```

## 参数

- Host
用于我们执行 SSH 命令的时候如何匹配到该配置。
  - *，匹配所有主机名。
  - *.example.com，匹配以 .example.com 结尾。
  - !*.dialup.example.com,*.example.com，以 ! 开头是排除的意思。
  - 192.168.0.?，匹配 192.168.0.[0-9] 的 IP。
- AddKeysToAgent
是否自动将 key 加入到 ssh-agent，值可以为 no(default)/confirm/ask/yes。

如果是 yes，key 和密码都将读取文件并以加入到 agent ，就像 ssh-add。其他分别是询问、确认、不加入的意思。添加到 ssh-agent 意味着将私钥和密码交给它管理，让它来进行身份认证。

- AddressFamily
指定连接的时候使用的地址族，值可以为 any(default)/inet(IPv4)/inet6(IPv6)。

- BindAddress
指定连接的时候使用的本地主机地址，只在系统有多个地址的时候有用。在 UsePrivilegedPort 值为 yes 的时候无效。

- ChallengeResponseAuthentication
是否响应支持的身份验证 chanllenge，yes(default)/no。

- Compression
是否压缩，值可以为 no(default)/yes。

- CompressionLevel
压缩等级，值可以为 1(fast)-9(slow)。6(default)，相当于 gzip。

- ConnectionAttempts
退出前尝试连接的次数，值必须为整数，1(default)。

- ConnectTimeout
连接 SSH 服务器超时时间，单位 s，默认系统 TCP 超时时间。

- ControlMaster
是否开启单一网络共享多个 session，值可以为 no(default)/yes/ask/auto。需要和 ControlPath 配合使用，当值为 yes 时，

ssh 会监听该路径下的 control socket，多个 session 会去连接该 socket，它们会尽可能的复用该网络连接而不是重新建立新的。

- ControlPath
指定 control socket 的路径，值可以直接指定也可以用一下参数代替：

- %L 本地主机名的第一个组件
- %l 本地主机名（包括域名）
- %h 远程主机名（命令行输入）
- %n 远程原始主机名
- %p 远程主机端口
- %r 远程登录用户名
- %u 本地 ssh 正在使用的用户名
- %i 本地 ssh 正在使用 uid
- %C 值为 %l%h%p%r 的 hash

请最大限度的保持 ControlPath 的唯一。至少包含 %h，%p，%r（或者 %C）。

- ControlPersist
结合 ControlMaster 使用，指定连接打开后后台保持的时间。值可以为 no/yes/整数，单位 s。如果为 no，最初的客户端关闭就关闭。如果 yes/0，无限期的，直到杀死或通过其它机制，如：ssh -O exit。

- GatewayPorts
指定是否允许远程主机连接到本地转发端口，值可以为 no(default)/yes。默认情况，ssh 为本地回环地址绑定了端口转发器。

- HostName
真实的主机名，默认值为命令行输入的值（允许 IP）。你也可以使用 %h，它将自动替换，只要替换后的地址是完整的就 ok。

- IdentitiesOnly
指定 ssh 只能使用配置文件指定的 identity 和 certificate 文件或通过 ssh 命令行通过身份验证，即使 ssh-agent 或 PKCS11Provider 提供了多个 identities。值可以为 no(default)/yes。

- IdentityFile
指定读取的认证文件路径，允许 DSA，ECDSA，Ed25519 或 RSA。值可以直接指定也可以用一下参数代替：

- %d，本地用户目录 ~
- %u，本地用户
- %l，本地主机名
- %h，远程主机名
- %r，远程用户名

- LocalCommand
指定在连接成功后，本地主机执行的命令（单纯的本地命令）。可使用 %d，%h，%l，%n，%p，%r，%u，%C 替换部分参数。只在 PermitLocalCommand 开启的情况下有效。

- LocalForward
指定本地主机的端口通过 ssh 转发到指定远程主机。格式：LocalForward [bind_address:]post host:hostport，支持 IPv6。

- PasswordAuthentication
是否使用密码进行身份验证，yes(default)/no。

- PermitLocalCommand
是否允许指定 LocalCommand，值可以为 no(default)/yes。

- Port
指定连接远程主机的哪个端口，22(default)。

- ProxyCommand
指定连接的服务器需要执行的命令。%h，%p，%r

如：ProxyCommand /usr/bin/nc -X connect -x 192.0.2.0:8080 %h %p

- User
登录用户名
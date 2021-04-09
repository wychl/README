# ansible

## 通过跳板机部署filebeat

### ssh配置文件

```config
Host filbeat-host
     HostName xxx.xxx.xxx.xxx
     User root
     Port 22
     ProxyCommand ssh -W %h:%p filbeat-bastion
     IdentityFile /path/to/filbeat-host/ssh-private-key

Host filbeat-bastion
     HostName xxx.xxx.xxx.xxx
     User root
     Port 22
     IdentityFile /path/to/filbeat-bastion/ssh-private-key

```

### ansible的host配置

```ansible-host
[filebeats]
filbeat-host ansible_ssh_common_args="-F /path/to/ssh-config"
```

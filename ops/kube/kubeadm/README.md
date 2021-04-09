# kubeadm

## 准备

### 关闭防火墙

如果各个主机启用了防火墙，需要开放Kubernetes各个组件所需要的端口，可以查看Installing kubeadm中的”Check required ports”一节。 这里简单起见在各节点禁用防火墙：

```sh
systemctl stop firewalld
systemctl disable firewalld
```

### 禁用SELINUX

```sh
# 临时禁用
setenforce 0

#永久禁用 
vim /etc/selinux/config    # 或者修改/etc/sysconfig/selinux
SELINUX=disabled
```

### 修改k8s.conf文件

```sh
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system
```

### 关闭swap

```sh
# 临时关闭
swapoff -a
```


修改 /etc/fstab 文件，注释掉 SWAP 的自动挂载（永久关闭swap，重启后生效）

```sh
# 注释掉以下字段
/dev/mapper/cl-swap     swap                    swap    defaults        0 0
```


## 安装Docker

### 卸载老版本的Docker

如果有没有老版本Docker，则不需要这步

```sh
yum remove docker \
           docker-common \
           docker-selinux \
           docker-engine
```

### 使用yum进行安装

```sh
# step 1: 安装必要的一些系统工具
sudo yum install -y yum-utils device-mapper-persistent-data lvm2
# Step 2: 添加软件源信息
sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
# Step 3: 更新并安装 Docker-CE
sudo yum makecache fast
sudo yum -y install docker-ce docker-ce-selinux
# 注意：
# 官方软件源默认启用了最新的软件，您可以通过编辑软件源的方式获取各个版本的软件包。例如官方并没有将测试版本的软件源置为可用，你可以通过以下方式开启。同理可以开启各种测试版本等。
# vim /etc/yum.repos.d/docker-ce.repo
#   将 [docker-ce-test] 下方的 enabled=0 修改为 enabled=1
#
# 安装指定版本的Docker-CE:
# Step 3.1: 查找Docker-CE的版本:
# yum list docker-ce.x86_64 --showduplicates | sort -r
#   Loading mirror speeds from cached hostfile
#   Loaded plugins: branch, fastestmirror, langpacks
#   docker-ce.x86_64            17.03.1.ce-1.el7.centos            docker-ce-stable
#   docker-ce.x86_64            17.03.1.ce-1.el7.centos            @docker-ce-stable
#   docker-ce.x86_64            17.03.0.ce-1.el7.centos            docker-ce-stable
#   Available Packages
# Step 3.2 : 安装指定版本的Docker-CE: (VERSION 例如上面的 17.03.0.ce.1-1.el7.centos)
sudo yum -y --setopt=obsoletes=0 install docker-ce-[VERSION] \
docker-ce-selinux-[VERSION]

# Step 4: 开启Docker服务
sudo systemctl enable docker && systemctl start docker
```

### docker安装问题小结

错误信息

```sh
Error: Package: docker-ce-17.03.2.ce-1.el7.centos.x86_64 (docker-ce-stable)
           Requires: docker-ce-selinux >= 17.03.2.ce-1.el7.centos
           Available: docker-ce-selinux-17.03.0.ce-1.el7.centos.noarch (docker-ce-stable)
               docker-ce-selinux = 17.03.0.ce-1.el7.centos
           Available: docker-ce-selinux-17.03.1.ce-1.el7.centos.noarch (docker-ce-stable)
               docker-ce-selinux = 17.03.1.ce-1.el7.centos
           Available: docker-ce-selinux-17.03.2.ce-1.el7.centos.noarch (docker-ce-stable)
               docker-ce-selinux = 17.03.2.ce-1.el7.centos
 You could try using --skip-broken to work around the problem
 You could try running: rpm -Va --nofiles --nodigest
```

解决办法


```sh
# 要先安装docker-ce-selinux-17.03.2.ce，否则安装docker-ce会报错
# 注意docker-ce-selinux的版本 要与docker的版本一直
yum install -y https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-selinux-17.03.2.ce-1.el7.centos.noarch.rpm

# 或者
yum -y install https://mirrors.aliyun.com/docker-ce/linux/centos/7/x86_64/stable/Packages/docker-ce-selinux-17.03.2.ce-1.el7.centos.noarch.rpm
```

参考链接

https://blog.csdn.net/csdn_duomaomao/article/details/79019764

### 安装校验

```sh
ocker version
Client:
 Version:      17.03.2-ce
 API version:  1.27
 Go version:   go1.7.5
 Git commit:   f5ec1e2
 Built:        Tue Jun 27 02:21:36 2017
 OS/Arch:      linux/amd64

Server:
 Version:      17.03.2-ce
 API version:  1.27 (minimum version 1.12)
 Go version:   go1.7.5
 Git commit:   f5ec1e2
 Built:        Tue Jun 27 02:21:36 2017
 OS/Arch:      linux/amd64
 Experimental: false
```


参考文献
https://yq.aliyun.com/articles/110806

## 安装kubeadm，kubelet，kubectl

在各节点安装kubeadm，kubelet，kubectl

### 修改yum安装源

```sh
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```

### 安装软件


```sh
yum install -y kubelet kubeadm kubectl
systemctl enable kubelet && systemctl start kubelet
```

## 初始化Master节点

### 配置kubeadm init 初始化文件

配置文件官网[介绍链接](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/?spm=a2c4e.10696291.0.0.7ed019a45gM2EF#config-file)

创建配置文件kubeadm-init.yaml文件

```yaml
apiVersion: kubeadm.k8s.io/v1beta2
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 0.0.0.0 #内网IP地址
  bindPort: 6443
nodeRegistration:
  criSocket: /var/run/dockershim.sock
  name: p1.morekoo.prd
  taints:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.cn-hangzhou.aliyuncs.com/google_containers
kind: ClusterConfiguration
kubernetesVersion: v1.16.0
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
scheduler: {}

```

### 运行初始化命令

```sh
kubeadm init --config kubeadm-init.yaml
```

### 使kubectl正常工作

```sh
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

## 安装Pod Network

```sh
kubectl apply -f https://docs.projectcalico.org/v3.8/manifests/calico.yaml
```

### master node参与工作负载

使用kubeadm初始化的集群，出于安全考虑Pod不会被调度到Master Node上，也就是说Master Node不参与工作负载。

这里搭建的是测试环境可以使用下面的命令使Master Node参与工作负载：

```sh
kubectl taint nodes node1 node-role.kubernetes.io/master-
```

## 向Kubernetes集群添加Node

### 命令

```sh
kubeadm join --token <token> <master-ip>:<master-port> --discovery-token-ca-cert-hash sha256:<hash>
```

### 查看token的值，在master节点运行以下命令
    
```sh
# 如果没有token，请使用命令kubeadm token create 创建
kubeadm token list
```

### 查看hash值，在master节点运行以下命令
    
```sh
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | \
   openssl dgst -sha256 -hex | sed 's/^.* //'
```

### node2加入集群很是顺利，下面在master节点上执行命令查看集群中的节点：
    
```sh
kubectl get nodes
```

### 如何从集群中移除Node

```sh
kubectl drain node2 --delete-local-data --force --ignore-daemonsets
kubectl delete node <node name>
```

在node2上执行：

```sh
kubeadm reset
```

## 安装ingress

[链接](https://github.com/nginxinc/kubernetes-ingress/blob/master/docs/installation.md)

## 安装dashboard

[链接](https://github.com/kubernetes/dashboard/blob/master/docs/user/installation.md)
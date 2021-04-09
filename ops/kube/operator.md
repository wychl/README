# operator 

## 安装`operator-sdk`

### 拉取operator-sdk

```sh
mkdir -p $GOPATH/src/github.com/operator-framework
cd $GOPATH/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
```
 
### 安装依赖（因为依赖包在墙外，需要设置代理）

```sh
http_proxy=http://x.x.x.x:port make dep
```

### 编译安装 `operator-sdk`

```sh
make install
```

## `operator`工作流

1. Create a new operator project using the SDK Command Line Interface(CLI)
2. Define new resource APIs by adding Custom Resource Definitions(CRD)
3. Define Controllers to watch and reconcile resources
4. Write the reconciling logic for your Controller using the SDK and controller-runtime APIs
6. Use the SDK CLI to build and generate the operator deployment manifests

## `operator`demo

### 创建名称为`app-operator`的项目

```sh
#创建工作目录
mkdir -p $GOPATH/src/github.com/example-inc/
# 在工作目录下创建app-operator项目
cd $GOPATH/src/github.com/example-inc/
operator-sdk new app-operator
cd app-operator
```

### 创建`AppService`自定义资源，版本为`app.example.com/v1alpha1`

```sh
operator-sdk add api --api-version=app.example.com/v1alpha1 --kind=AppService
```

### 创建`AppService`自定义资源的`controller`

```sh
operator-sdk add controller --api-version=app.example.com/v1alpha1 --kind=AppService
```

### 构建`operator`镜像，并将镜像推到`registry`

将 `wanyanchengli/app-operator:v1`替换成自己的registry地址

```sh
operator-sdk build wanyanchengli/app-operator:v1
docker push wanyanchengli/app-operator:v1
```

### 将`operator manifest`的镜像更新为`wanyanchengli/app-operator:v1`

```sh
# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)
sed -i 's|REPLACE_IMAGE|quay.io/example/app-operator|g' deploy/operator.yaml
# Mac系统
sed -i "" 's|REPLACE_IMAGE|quay.io/example/app-operator|g' deploy/operator.yaml
```

### 设置`operator`的权限

1. 创建`service account` 

```sh
kubectl create -f deploy/service_account.yaml
```

2. 设置`RBAC`

```sh
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
```

### 在k8s中创建CRD 

```sh
kubectl create -f deploy/crds/app_v1alpha1_appservice_crd.yaml
```

### 部署`operator` 

```sh
kubectl create -f deploy/operator.yaml
```

### Create an AppService CR

```sh
### The default controller will watch for AppService objects and create a pod for each CR
kubectl create -f deploy/crds/app_v1alpha1_appservice_cr.yaml
```

### 验证`operator`是否创建成功

```sh
kubectl get pod -l app=example-appservice
# NAME                     READY     STATUS    RESTARTS   AGE
# example-appservice-pod   1/1       Running   0          1m
```

### 清理`operator`

```sh
kubectl delete -f deploy/crds/app_v1alpha1_appservice_cr.yaml
kubectl delete -f deploy/operator.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/crds/app_v1alpha1_appservice_crd.yaml
```
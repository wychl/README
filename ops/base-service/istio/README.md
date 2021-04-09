# istio

## helm安装istio

- 下载包

```sh
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.1.2 sh -
cd istio-1.1.2
```

- 添加helm仓库

```sh
helm repo add istio.io https://storage.googleapis.com/istio-release/releases/1.1.2/charts/
```

- 使用`helm template`安装
  - 创建`istio-system` namespace

    ```sh
    kubectl create namespace istio-system
    ```

  - 注册istio CRD

    ```sh
    helm template install/kubernetes/helm/istio-init --name istio-init --namespace istio-system | kubectl apply -f -
    ```
  
  - 验证53个CRD是否在k8s中注册成功(cert-mamaget启用时58个)

    ```sh
    kubectl get crds | grep 'istio.io\|certmanager.k8s.io' | wc -l
    ```
  
  - 安装

    ```sh
    helm template install/kubernetes/helm/istio --name istio --namespace istio-system | kubectl apply -f -
    ```

- 验证istio是否安装成功

```sh
kubectl get svc -n istio-system
kubectl get pods -n istio-system
```

## 将`istio-demo` namespace 设置为自动sidecar注入

```sh
#创建namespace
kubectl create ns istio-demo

#设置namespace label
kubectl label namespace istio-demo istio-injection=enabled

# 检测
kubectl get namespace -L istio-injection
#NAME           STATUS   AGE    ISTIO-INJECTION
#istio-demo     Active   29s    enabled

```
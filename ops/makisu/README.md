# makisu

## 创建 registry 配置

```sh
kubectl create secret generic docker-registry-config --from-file=./registry.yaml
```

## 构建镜像-推送到registry

```sh
kubectl create -f job.yaml
```

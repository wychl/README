# 使用多阶段构建镜像减少镜像文件大小

## 多阶段构构建镜像

```sh
docker build -t dev-to:multi-stage .
# 查看镜像大小
docker images | grep dev-to -B 1
```

### 结果

***镜像大小：7.18MB***

## 单段构构建镜像

```sh
docker build -t dev-to:single-stage . -f single-stage-build
# 查看镜像大小
docker images | grep dev-to -B 1
```

### 结果

***镜像大小：404MB***

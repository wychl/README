# k8s 笔记

## configMap使用

### 创建configMap

- 语法 `kubectl create configmap <map-name> <data-source>`
  - `map-name`:创建configMap时指定的名称
  - `data-source`:目录、文件或者字面值

### 示例

- 使用字面值创建ConfigMap `kubectl create configmap <map-name> --from-literal=key=value --from-literal=key2=value2`

```sh
kubectl create configmap demo-config --from-literal=demo.port=8888
```

-询创建的ConfigMap

```sh
kubectl get configmap demo-config -o yaml
```

返回的结果

``` yaml
apiVersion: v1
data:
  demo.port: "8888"
kind: ConfigMap
metadata:
  creationTimestamp: 2018-12-04T13:44:02Z
  name: demo-config
  namespace: default
  resourceVersion: "66775"
  selfLink: /api/v1/namespaces/default/configmaps/demo-config
  uid: a8f73ecd-f7ca-11e8-ad3e-080027dedea1
```

- 使用ConfigMap
  - 以环境变量的方式使用ConfigMap

pod规格文件

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-demo
spec:
  containers:
    - name: test-demo-container
      image: demo:1.0
      volumeMounts:
      - name: config-volume
        mountPath: /etc/config #ConfigMap数据存放位置
      env:
      # 定义环境变量
        - name: PORT
          valueFrom:
            configMapKeyRef:
              # ConfigMap名称
              name: demo-config
              # 引用的key
              key: demo.port
  volumes:
    - name: config-volume
      configMap:
        # ConfigMap名称
        name: demo-config
  restartPolicy: Never
```

  - 将ConfigMap数据加入到Volume
  
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-demo
spec:
  containers:
    - name: test-demo-container
      image: demo:1.0
      volumeMounts:
      - name: config-volume
        mountPath: /etc/config #ConfigMap数据存放位置
  volumes:
    - name: config-volume
      configMap:
        # ConfigMap名称
        name: demo-config
  restartPolicy: Never
```

```sh
# 打包镜像
./build.sh

#将镜像同步到minikube
./sysnc_minikube.sh

# 创建pod
kubectl create -f demo.yaml

#进入pod中
 kubectl exec -it test-demo sh

 # 验证结果
 echo $PORT # 8888
    

 ls /etc/config/ #demo.port
```

#### 使用文件创建ConfigMap

- 创建ConfigMap
  
```bash
kubectl create configmap example-demo-file-config --from-file=demo-config.json
```

结果

```yaml
apiVersion: v1
data:
  demo-config.json: |-
    {
        "debug":true,
        "port":8080
    }
kind: ConfigMap
metadata:
  creationTimestamp: 2018-12-04T15:49:13Z
  name: example-demo-file-config
  namespace: default
  resourceVersion: "75298"
  selfLink: /api/v1/namespaces/default/configmaps/example-demo-file-config
  uid: 2610525f-f7dc-11e8-be9e-080027dedea1
```

- 使用ConfigMap

pod规格文件

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-demo
spec:
  containers:
    - name: test-demo-container
      image: demo:1.0
      volumeMounts:
      - name: config-volume
        mountPath: /etc/config #ConfigMap数据存放位置
  volumes:
    - name: config-volume
      configMap:
        # ConfigMap名称
        name: example-demo-file-config
        items:
        - key: demo-config.json
          path: config.json
  restartPolicy: Never
```  

- 测试

```sh
# 打包镜像

# 创建pod
kubectl create -f demo.yaml

#进入pod中
 kubectl exec -it test-demo sh

 # 验证结果    

 ls /etc/config/ #demo.port
```

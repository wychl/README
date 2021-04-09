# TLS 使用说明

## 介绍

TLS(Transport Layer Security)为通过网络进行通信的应用程序的传输的数据加密。HTTPS(Hypertext Transfer Protocol Secure)
是HTTP的拓展，利用TLSt提高安全性。TLS技术要求CA(Certificate Authority)向服务提供方颁发X.509数字证书，并将证书移交给服务消费者便于CA验证。

## 生成服务端证书

### 方式1

1. 生成私钥

```sh
openssl genrsa -out server.key 2048
```

2. 生成证书

```sh
openssl req -new -x509 -key server.key -out server.pem -days 3650
```

### 方式2

生成2048位bit密钥，基于密钥生成10年有效期的证书，此外，CN = localhost声明该证书对localhost域有效。

```sh
openssl req -newkey rsa:2048 \
  -new -nodes -x509 \
  -days 3650 \
  -out server.pem \
  -keyout server.key \
  -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost"
```

## 使用案例

### 单向认证-客户端没有服务端证书

- 启动服务端(需要提前生成服务端密钥和证书文件)

```sh
go run go-tls-single-demo/server/main.go
```

- 客户端没有设置跳过认证

```sh
go run go-tls-single-demo/client/main.go
```

证书验证出错,输出结果为

    ```
    2019/05/29 22:22:38 Get https://localhost:8443/hello: x509: certificate signed by unknown authority
    exit status 1
    ```

- 客户端设置跳过认证

```sh
go run go-tls-single-demo/client-skip-verify/main.go
```

客户端告诉服务端不进行证书验证，输出结果为

    ```
    Hello, world!
    ```





### 单向认证-客户端使用服务端证书

- 启动服务端(服务端和客户端公用一个证书)

```sh
go run go-tls-single-demo/server/main.go
```

- 客户端没有设置跳过认证

```sh
go run go-tls-single-demo/client/main.go
```

证书验证出错,输出结果为

    ```
    2019/05/29 22:22:38 Get https://localhost:8443/hello: x509: certificate signed by unknown authority
    exit status 1
    ```

- 客户端设置跳过认证

```sh
go run go-tls-single-demo/client-server-share-crt/main.go
```

客户端告诉服务端不进行证书验证，输出结果为

    ```
    Hello, world!
    ```

### 双向认证-分别生成服务端和客户端证书

服务端证书生成方式保持不变，而客户端证书生成方式如下。

1. 生成私钥

```sh
openssl genrsa -out client.key 2048
```

2. 为客户端创建中间证书

```sh
openssl req -new -key  client.key -out client.csr \
  -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost"
```

3. 生成证书
```sh
openssl x509 -req -days 365 -in client.csr -CA server.pem -CAkey server.key -CAcreateserial -out client.pem
```

- 启动服务端(需要提前生成服务端密钥和证书文件)

```sh
go run go-tls-manual-demo/server/main.go
```

- 通过客户端进行测试

```sh
go run go-tls-manual-demo/client/main.go
```

证书验证输出结果为

    ```
    Hello, world!
    ```

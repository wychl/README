# mkcert

## 生成证书

```sh
mkcert -install
 mkcert -cert-file certs/domain.cert -key-file certs/domain.key localhost
```

## 测试

```sh
docker run -d -p 8080:80 -p 8443:443 \
-v ${PWD}/default.conf:/etc/nginx/conf.d/default.conf \
-v ${PWD}/certs/domain.cert:/root/domain.cert \
-v ${PWD}/certs/domain.key:/root/domain.key \
nginx:alpine
```

地址 https://localhost:8443/
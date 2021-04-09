# rest api


## 测试

- allow

```sh
curl -H 'user:alice' http://127.0.0.1:8080/users/1
```

- deny

```sh
curl -H 'user:alice' http://127.0.0.1:8080/users/3
```

- 超级用户

```sh
curl -H 'user:super' http://127.0.0.1:8080/users/3
```
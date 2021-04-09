# filebeat

## filebeat与elasticsearch的链接方式

### `baic authentication`

```yaml
output.elasticsearch:
  username: filebeat 
  password: verysecret 
  protocol: https 
  hosts: ["elasticsearch.example.com:9200"] 
```

- `username`: elasticsearch账号
- `password`: elasticsearch密码
- `protocol`: 链接协议
- `hosts`: elasticsearch地址数组

### 证书的方式

```yaml
output.elasticsearch:
  username: filebeat
  password: verysecret
  protocol: https
  hosts: ["elasticsearch.example.com:9200"]
  ssl.certificate_authorities: 
    - /etc/pki/my_root_ca.pem
    - /etc/pki/my_other_ca.pem
  ssl.certificate: "/etc/pki/client.pem" 
  ssl.key: "/etc/pki/key.pem" 
```

- `ssl.certificate_authorities` CA证书数组
- `ssl.certificate` 客户端证书
- `ssl.key` 客户端证书密钥

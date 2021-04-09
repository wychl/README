# vault

## 什么是Vault

HashiCorp Vault是一款企业级私密信息管理工具。HashiCorp是一家专注于DevOps工具链的公司，其旗下明星级产品包括Vagrant、Packer、Terraform、Consul、Nomad等，再加上Vault，这些工具贯穿了持续交付的整个流程。

## API操作

### 初始化

命令

```sh
curl -X PUT -d "{\"secret_shares\":1, \"secret_threshold\":1}"  http://127.0.0.1:8200/v1/sys/init | jq
```

输出

```json
{
  "keys": [
    "5a94a35f42813154076aafe6a4884e2218a140873fb006ae6b253bc976fe25e5"
  ],
  "keys_base64": [
    "WpSjX0KBMVQHaq/mpIhOIhihQIc/sAauayU7yXb+JeU="
  ],
  "root_token": "s.dyXbo3SSCTg07DFqd8Z6UfF9"
}
```

第一个是master key的public key，第二个是unseal key,最后一个是root token。unseal vault之后才能验证进行具体的操作。

### unseal

命令

```sh
curl -X PUT -d '{"key": "WpSjX0KBMVQHaq/mpIhOIhihQIc/sAauayU7yXb+JeU="}'  http://127.0.0.1:8200/v1/sys/unseal | jq

```

输出

```json
{
  "type": "shamir",
  "initialized": true,
  "sealed": false,
  "t": 1,
  "n": 1,
  "progress": 0,
  "nonce": "",
  "version": "1.1.3",
  "migration": false,
  "cluster_name": "vault-cluster-983e319d",
  "cluster_id": "5a2b8523-d479-0993-e03f-29a11b631220",
  "recovery_seal": false
}
```

### 创建新token

为了安全起见，我们可以用root token创建出有限权限的新token，来继续后面的操作。假设这个token可以读写的路径为secret/*。在这个之前我们先创建访问这个路径的policy:

- 创建admin policy

```sh
# 创建 admin policy
curl -v  -X PUT  -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9" -d '{"policy": "path \"secret/*\" {\n  capabilities = [\"create\"] \n}"}'   http://127.0.0.1:8200/v1/sys/policy/admin-policy


# 查询 admin policy

curl -X GET  -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9"  http://127.0.0.1:8200/v1/sys/policy/admin-policy | jq

```

admin policy

```json
{
  "rules": "path \"secret/*\" {\n  capabilities = [\"create\"] \n}",
  "name": "admin-policy",
  "request_id": "738d9523-6645-223e-88c3-4a10fad42dd8",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "name": "admin-policy",
    "rules": "path \"secret/*\" {\n  capabilities = [\"create\"] \n}"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

- 创建 user policy

```sh
# 创建 admin policy
curl -v  -X PUT  -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9" -d '{"policy": "path \"secret/*\" {\n  capabilities = [\"read\"] \n}"}'   http://127.0.0.1:8200/v1/sys/policy/user-policy


# 查询 admin policy

curl -X GET  -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9"  http://127.0.0.1:8200/v1/sys/policy/user-policy | jq
```

user policy

```json
{
  "name": "user-policy",
  "rules": "path \"secret/*\" {\n  capabilities = [\"read\"] \n}",
  "request_id": "4aa7f0eb-fec2-c243-b263-50f74ca79a5d",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "name": "user-policy",
    "rules": "path \"secret/*\" {\n  capabilities = [\"read\"] \n}"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

- 使用admin policy创建token

命令

```sh
curl -X POST -d '{"policies": ["admin-policy"]}' -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9"  http://127.0.0.1:8200/v1/auth/token/create | jq
```

输出

```json
{
  "request_id": "c43674f2-6b2a-e814-4095-01583b545068",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "s.ppKxTh1WwSbagWYrKojacMEm",
    "accessor": "vHIZZfrhxoda0mk0W9E6jNP4",
    "policies": [
      "admin-policy",
      "default"
    ],
    "token_policies": [
      "admin-policy",
      "default"
    ],
    "metadata": null,
    "lease_duration": 604800,
    "renewable": true,
    "entity_id": "",
    "token_type": "service",
    "orphan": false
  }
}
```

- 使用user policy创建token


命令

```sh
curl -X POST -d '{"policies": ["user-policy"]}' -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9"  http://127.0.0.1:8200/v1/auth/token/create | jq
```

输出

```json
{
  "request_id": "d6898b1f-ed4e-ce10-563b-4bef9becad56",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "s.km3x394gwEPJxpRBvleoUofx",
    "accessor": "RR8oTZ43y8RxhPEriauWAXcr",
    "policies": [
      "default",
      "user-policy"
    ],
    "token_policies": [
      "default",
      "user-policy"
    ],
    "metadata": null,
    "lease_duration": 604800,
    "renewable": true,
    "entity_id": "",
    "token_type": "service",
    "orphan": false
  }
}
```

- 测试secrets/*路径下的token权限

```sh
export ADMIN_TOKEN="s.dyXbo3SSCTg07DFqd8Z6UfF9"
curl -X POST -H "X-Vault-Token:s.dyXbo3SSCTg07DFqd8Z6UfF9" -d '{"token":"c192d0211cb81fbfeee53fb16e2a7465"}' http://127.0.0.1:8200/v1/secret/config



export USER_TOKEN="s.km3x394gwEPJxpRBvleoUofx"
curl -X GET -H "X-Vault-Token:$USER_TOKEN" http://127.0.0.1:8200/v1/secret/hello | jq

```

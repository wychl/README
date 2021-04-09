# opa

## demo测试

- 启动测试环境

```sh
docker-compose -f docker-compose.yml up -d
```


- 将示例策略加载的opa

```sh
curl -X PUT --data-binary @example.rego \
  localhost:8181/v1/policies/example
```  

- 查看alice是否能查看自己的工资

```sh
curl --user alice:password localhost:5000/finance/salary/alice
```

-  bob是否能查看alice的工资(因为bob是alice的上级)

```sh
curl --user bob:password localhost:5000/finance/salary/alice
```

-  bob是否能查看 charlie的工资（bob不是charlie的上级）.

```sh
curl --user bob:password localhost:5000/finance/salary/charlie
```

- 改变策略（允许HR查看所有人的工资）


cat >example-hr.rego <<EOF
package httpapi.authz

import input

# 允许HR查看所有人的工资

将策略写入到opa

```
curl -X PUT --data-binary @example-hr.rego \
  http://localhost:8181/v1/policies/example-hr
```


检验修改后的策略

```sh
curl --user david:password localhost:5000/finance/salary/alice
curl --user david:password localhost:5000/finance/salary/bob
curl --user david:password localhost:5000/finance/salary/charlie
curl --user david:password localhost:5000/finance/salary/david
```
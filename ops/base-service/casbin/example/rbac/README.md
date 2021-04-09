# RBAC

## policy

```csv
g, alice, data2_admin
```

**解释：**

alice是data2_admin的成员，alice可以理解为一个人，一个资源或者一个角色等任何东西，但是在casbin眼里就是一个字符串没有任何意义。

## matcher

```conf
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

## pattern matching

```csv
p, alice, book_group, read
g, /book/1, book_group
g, /book/2, book_group
```

alice可以读取`book_group`中的书(`/book/1`,`/book/2`)

```csv
g, /book/:id, book_group
```
可以匹配任何`/book/:int`格式的书


## 注意事项

1. user和role不能相同，因为casbin将user和role解释为字符串 user_alice role_alice
2. 角色继承，比如alice有role1，role1有role2，则alice有role，cashbin可以通过`NewRoleManager(maxHierarchyLevel int)`设置等级


## RBAC with Domains(多租户)

```conf
[role_definition]
g = _, _, _
```

```csv
p, admin, tenant1, data1, read
p, admin, tenant2, data2, read

g, alice, admin, tenant1
g, alice, user, tenant2
```
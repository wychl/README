# 文档

model file 定义规则, policy file 中定义的就是具体的内容.


## Casbin做了什么
支持自定义请求的格式，默认的请求格式为{subject, object, action}。
具有访问控制模型model和策略policy两个核心概念。
支持RBAC中的多层角色继承，不止主体可以有角色，资源也可以具有角色。
支持超级用户，如root或Administrator，超级用户可以不受授权策略的约束访问任意资源。
支持多种内置的操作符，如keyMatch，方便对路径式的资源进行管理，如/foo/bar可以映射到/foo*

## Casbin不做的事情
身份认证authentication(即验证用户的用户名、密码)，casbin只负责访问控制。应该有其他专门的组件负责身份认证，然后由casbin进行访问 控制，二者是相互配合的关系。
管理用户列表或角色列表。Casbin认为由项目自身来管理用户、角色列表更为合适，用户通常有他们的密码，但是Casbin的设计思想并不是把 它作为一个存储密码的容器。而是存储RBAC方案中用户和角色之间的映射关系。

## PERM模型

PERM(Policy, Effect, Request, Matchers)模型很简单, 但是反映了权限的本质 – 访问控制

- Policy: 定义权限的规则
- Effect: 定义组合了多个 Policy 之后的结果, allow/deny
- Request: 访问请求, 也就是谁想操作什么
- Matcher: 判断 Request 是否满足 Policy

![perm](./img/perm.png)

## model语法

至少包含 `[request_definition]`, `[policy_definition]`, `[policy_effect]`, `[matchers]`四部分;RBAC model还需要`[role_definition]`

### Request definition
 
定义访问请求，即`e.Enforce(...)`函数的参数;`sub`定义访问实体,`obj`访问资源,`act`访问方法

### Policy定义

- model

```conf
[policy_definition]
p = sub, obj, act
p2 = sub, act
```

p,p2策略规则名称

- role定义

```conf
[role_definition]
g = _, _
g2 = _, _
g3 = _, _, _
```


g, g2, g3 表示不同的RBAC 体系;_, _ 表示用户和角色; _, _, _ 表示用户, 角色, 域(也就是租户)

- policy(csv格式)
 
每一行对应一个规则

```csv
p, alice, data1, read
p2, bob, write-all-objects
```

- 对应关系

```
(alice, data1, read) -> (p.sub, p.obj, p.act)
(bob, write-all-objects) -> (p2.sub, p2.act)
```

## Policy effect

当多个策略规则应用于访问请求时最终结果该如何；比如一个规则允许，另一个规则禁止境况。

```conf
[policy_effect]
e = some(where (p.eft == allow))
```

- 上述`policy effect`表明：只要任何一个策略规则是`allow`，则最终的`effect`是`allow`
- `p.eft`的值是`allow`或者`deny`,默认值为`allow`，所以上面为指定


```conf
[policy_effect]
e = !some(where (p.eft == deny))
```

- 如果没有匹配结果为`deny`，则最终的`effect`为`allow`
- `some`,存在匹配规则，`any`所有的匹配规则

### matcher

定义策略表达式

```csv
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```

自定义request和policy匹配方式， p.eft 是 allow 还是 deny, 就是基于此来决定的。

```go
func KeyMatchFunc(args ...interface{}) (interface{}, error) {
    name1 := args[0].(string)
    name2 := args[1].(string)

    return (bool)(KeyMatch(name1, name2)), nil
}

e.AddFunction("my_func", KeyMatchFunc)
```

使用自定义`key match`

```conf
[matchers]
m = r.sub == p.sub && my_func(r.obj, p.obj) && r.act == p.act
```


## model存储

- 文件加载


```go
e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
```

- 代码加载

```go
m := casbin.NewModel()
m.AddDef("r", "r", "sub, obj, act")
m.AddDef("p", "p", "sub, obj, act")
m.AddDef("g", "g", "_, _")
m.AddDef("e", "e", "some(where (p.eft == allow))")
m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

// Load the policy rules from the .CSV file adapter.
// Replace it with your adapter to avoid files.
a := persist.NewFileAdapter("examples/rbac_policy.csv")

// Create the enforcer.
e := casbin.NewEnforcer(m, a)
```

- 字符串加载

```go
// Initialize the model from a string.
text :=
`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
m := NewModel(text)

// Load the policy rules from the .CSV file adapter.
// Replace it with your adapter to avoid files.
a := persist.NewFileAdapter("examples/rbac_policy.csv")

// Create the enforcer.
e := casbin.NewEnforcer(m, a)
```

## policy存储
- 文件加载

```go
import "github.com/casbin/casbin"

e := casbin.NewEnforcer("examples/basic_model.conf", "examples/basic_policy.csv")
```

or 

```go
import (
    "github.com/casbin/casbin"
    "github.com/casbin/casbin/file-adapter"
)

a := fileadapter.NewAdapter("examples/basic_policy.csv")
e := casbin.NewEnforcer("examples/basic_model.conf", a)
```

- MySQL

```go
import (
    "github.com/casbin/casbin"
    "github.com/casbin/mysql-adapter"
)

a := mysqladapter.NewAdapter("mysql", "root:@tcp(127.0.0.1:3306)/")
e := casbin.NewEnforcer("examples/basic_model.conf", a)
```

- 自定义存储

  - 实现一个存储(实现四个接口)

    ```go
    // AddPolicy adds a policy rule to the storage.
    func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
        return errors.New("not implemented")
    }

    // RemovePolicy removes a policy rule from the storage.
    func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
        return errors.New("not implemented")
    }

    // RemoveFilteredPolicy removes policy rules that match the filter from the storage.
    func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
        return errors.New("not implemented")
    }
    ```  

  - 使用自定义存储

    ```go
    import (
        "github.com/casbin/casbin"
        "github.com/your-username/your-repo"
    )

    a := yourpackage.NewAdapter(params)
    e := casbin.NewEnforcer("examples/basic_model.conf", a)
    ```

- 运行时加载model或者加载、保存policy

```go
// 从 CONF 文件加载model
e.LoadModel()

// 从存储加载policy
e.LoadPolicy()

// 将当前policy保存到存储空间
e.SavePolicy()

```

- AutoSave

```go
import (
    "github.com/casbin/casbin"
    "github.com/casbin/xorm-adapter"
    _ "github.com/go-sql-driver/mysql"
)

// By default, the AutoSave option is enabled for an enforcer.
a := xormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
e := casbin.NewEnforcer("examples/basic_model.conf", a)

// 禁用AutoSave
e.EnableAutoSave(false)

// AutoSave禁止的，所以只影响当前enforer
// 不影响存储的Policy
e.AddPolicy(...)
e.RemovePolicy(...)

// 启用Autosave
e.EnableAutoSave(true)

// AutoSave启用的，不仅影响当前enforcer，也影响存储的policy
e.AddPolicy(...)
e.RemovePolicy(...)
```


## rest api


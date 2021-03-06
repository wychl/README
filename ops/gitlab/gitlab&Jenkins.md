# gitlab和jenkins集成

## 创建`Gitlab`账户

`Jenkins`使用此账号clone仓库代码

***步骤如下：***

### 第一步进入配置页面

点击icon进入配置页面

![](images/gitlab_config_entry.png)

### 第二步进入用户新建页面

现在配置页面首先点击`用户`->然后点击`新用户`，进入用户新建页面

![](images/gitlab_config.png)

### 第三步

用户新建页面配置账号

![](images/gitlab_user_form.png)

### 激活账号

点击邮箱链接激活

## 配置`Jenkins`凭证

***使用上步骤中账号和密码创建凭证***

### 第一步进入配置页面

![](images/genkins_admin_config.png)

#### 第二步

点击`凭据`->点击`.global`进入凭据加入页面

![](images/jenkins_global_config.png)

### 第三步

点击`Add Credentials`

![](images/add_credentials.png)

### 配置 `Credentials`

![](images/add_credentials_form.png)

需要编辑的参数如下

- Username `xxx` gitlab账号
- Password `xxxx` gitlab密码
- ID `xxx` 唯一
- Description `xxx`

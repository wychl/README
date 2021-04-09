# 全局忽略

1. 创建 `.gitignore_global`文件： `touch ~/.gitignore_global`
2. 配置git全局忽略 `git config --global core.excludesfile ~/.gitignore_global`
3. 查看git配置 `git config --list`

## 设置git编辑器

`git config --global core.editor vim`

## 删除文件历史

1. To remove the file, enter git rm --cached:

```sh
$ git rm --cached giant_file
# Stage our giant file for removal, but leave it on disk
```

2. Commit this change using --amend -CHEAD:

```sh
$ git commit --amend -CHEAD
# Amend the previous commit with your change
# Simply making a new commit won't work, as you need
# to remove the file from the unpushed history as well
```

3. Push your commits to GitHub:

```sh
$ git push
# Push our rewritten, smaller commit
```

## 使用代码操作`git server`

```sh
go run main.go \
    -user=git_user \
    -password=git_password \
    -repo=https://github.com/wychl/helloworld.git \
    -dir=./helloworld
```

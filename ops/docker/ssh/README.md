# SSH服务镜像

- [Dockerfile](./Dockerfile)
    - port 22
    - passwd screencast
    - user root
- 构建镜像 `docker build -t ssh .`
- 启动ssh服务 `docker run -p 222:22 -d --name ssh ssh:latest`
- 测试 `ssh root@127.0.0.1 -p222`
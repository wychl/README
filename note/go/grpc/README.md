# grpc

## 安装protoc

- linux

```sh
mkdir temp && cd temp
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.0/protoc-3.9.0-linux-x86_64.zip
unzip protoc-3.9.0-linux-x86_64.zip
sudo cp bin/protoc /usr/bin/
sudo cp -r include/google/protobuf/ /usr/include/google/
chmod 755 /usr/bin/protoc
sudo chmod -R 755 /usr/include/google/protobuf
```

- mac

```sh
mkdir temp && cd temp
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.0/protoc-3.9.0-osx-x86_64.zip
unzip protoc-3.9.0-osx-x86_64.zip
sudo cp bin/protoc /usr/local/bin/
chmod 755 /usr/local//bin/protoc
sudo cp -r include/google/protobuf/ /usr/local/include/google/
sudo chmod -R 755 /usr/local/include/google/
```


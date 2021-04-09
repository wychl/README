export GOOS=linux
go build -o server
docker build -t wanyanchengli/processer:v1 .
docker push wanyanchengli/processer:v1

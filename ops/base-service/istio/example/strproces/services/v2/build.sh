export GOOS=linux
go build -o server
docker build -t wanyanchengli/processer:v2 .
docker push wanyanchengli/processer:v2

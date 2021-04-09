export GOOS=linux
go build -o server
docker build -t wanyanchengli/processer:v3 .
docker push wanyanchengli/processer:v3

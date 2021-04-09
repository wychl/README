export GOOS=linux
go build -o server
docker build -t wanyanchengli/ui:v1 .
docker push wanyanchengli/ui:v1

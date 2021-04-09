docker build -t alertmanager-webhook:v1 .
docker run --name webhook -p 5001:5001 -d alertmanager-webhook:v1
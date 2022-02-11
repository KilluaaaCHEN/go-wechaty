git pull
name="go-wechaty"
version="1.0"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./runtime/main ./main.go
docker build -t "${name}:${version}" -f "./Dockerfile" .
docker stop ${name}
docker rm ${name}
docker run -d --name="${name}" ${name}:${version}
docker logs -f ${name}
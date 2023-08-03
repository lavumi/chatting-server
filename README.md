# HTTP2 GameServer

Http2 game-server made by gin



## Project Layout

Refer to [golang-standards/project-layout](https://github.com/golang-standards/project-layout)



## Run Server
```bash
go run cmd/server.go
```


## Build
todo improve
```bash
env GOOS=linux go build -o game-server cmd/server.go
docker build -t game-server:0.0.2 .
docker save -o game-server.tar game-server
rm game-server
```
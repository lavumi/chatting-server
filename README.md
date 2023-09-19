# HTTP2 Chatting Server

Http2 chat-server made by gin



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

### uglify js

```bash
uglifyjs web/js/netHandler.js web/js/script.js -o temp_build/web/js/chat.min.js -c -m
```

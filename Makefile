all: 
	go build -o docker/poller docker/poller.go 
	go build -o webserver/server webserver/server.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker/poller_linux docker/poller.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webserver/server_linux webserver/server.go

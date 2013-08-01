all: 
	go build -o poller/poller poller/poller.go 
	go build -o web web.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o poller/poller_linux poller/poller.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web web.go

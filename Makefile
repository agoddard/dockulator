all: 
	go build -o poller/poller poller/poller.go 
	go build -o web web.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o poller/poller_linux poller/poller.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web web.go

docker: centos ubuntu

centos: 
	# This has to have a ;\ at the end because we have to be in that dir to execute build
	cp -R calculators dockerfiles/centos/
	cd dockerfiles/centos;\
	docker build -t "dockulator-centos" .

ubuntu: 
	# This has to have a ;\ at the end because we have to be in that dir to execute build
	cp -R calculators dockerfiles/ubuntu/
	cd dockerfiles/ubuntu;\
	docker build -t "dockulator-ubuntu" .

clean:
	rm -rf dockerfiles/ubuntu/calculators
	rm -rf dockerfiles/ubuntu/calculators
	rm -f poller/poller_linux
	rm -f poller/poller
	rm -f web

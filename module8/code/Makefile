.PHONY : build
build :
	GOOS=linux GOARCH=amd64 go build  -o main . 

.PHONY : image
image :
	docker build -t httpserver:v1 . 

.PHONY : rundocker
rundocker :
	docker run -p 8081:8083 --env GOPORT=8083 --env GOAPPPATH=/data/logs/ -v /Users/darren/go/src/github.com/darrenli6/CloudNativeGo/module8/code/log/:/data/logs/  --name httpserver -itd httpserver:v1

.PHONY : killdocker
killdocker :
	docker rm -f httpserver






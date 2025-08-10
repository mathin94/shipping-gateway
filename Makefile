.PHONY: build build-server run

build:
	go build -o shipping-gateway

build-server:
	GOOS=linux GOARCH=amd64 go build -o shipping-gateway

run:
	go run cmd/web/main.go
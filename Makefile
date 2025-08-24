.PHONY: build build-server run

build:
	go build -o shipping-aggregator cmd/web/main.go

build-server:
	GOOS=linux GOARCH=amd64 go build -o shipping-aggregator cmd/web/main.go

run:
	go run cmd/web/main.go
	
tidy:
	go mod tidy
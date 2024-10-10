build: 
	@go build -o api

run: build
	@./api

test: 
	@go test ./...

cover:
	@go test -v -race -coverprofile coverage.txt ./...
	@go tool cover -html coverage.txt -o cover.html

lint:
	@golangci-lint run ./...

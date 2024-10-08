build: 
	@go build -o api

run: build
	@./api

test: 
	@go test ./...

lint:
	@golangci-lint run ./...

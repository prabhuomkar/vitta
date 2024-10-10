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

migrate-up:
	@migrate -path=database/migrations -database "${DATABASE_URL}" -verbose up

migrate-down:
	@migrate -path=database/migrations -database "${DATABASE_URL}" -verbose down

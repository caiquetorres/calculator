build:
	@go build -o bin/calculator cmd/main.go

run: build
	@./bin/calculator

test:
	@go test ./... -v

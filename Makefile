build:
	@go build -o bin/encoder cmd/main.go

run: build
	@./bin/encoder

test:
	@go test ./... -v

build:
	@go build -o bin/drain-node

run: build 
	@./bin/drain-node

test:
	go test ./... -v
build:
	@go build -o bin/databasetester

run: build
	@./bin/databasetester

test:
	@go test -v ./...
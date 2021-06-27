all: lint test

lint:
	golangci-lint run ./...

test:
	go test -cover -count=1 ./...

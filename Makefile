default: check test
test:
	go test -v -race -count=1 ./...
check:
	golangci-lint run -v ./...
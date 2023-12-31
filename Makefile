.PHONY: test
test:
	go test -v -race ./...

.PHONY: test/cover
test/cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

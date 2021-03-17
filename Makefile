
build: dist/go-oauth-curl

dist/go-oauth-curl:
	go build -o $@ ./...

test:
	golangci-lint run ./...

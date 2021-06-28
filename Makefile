
build: dist/beareq

dist/beareq:
	go build -o $@ ./...

update:
	go get -u -d -v ./...
	go mod tidy

test:
	go vet ./...
	go test -coverprofile cover.out -v ./...
	golangci-lint run --fix ./...

add_test:
	git ls-files | grep '.*\.go$$' | grep -v "_test\.go$$" | while read -r f; do gotests -all $$f | grep -v "^Generated" > $${f%%.go}_test.go; done

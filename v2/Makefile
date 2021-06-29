install:
	go install -v ./cmd/...

update:
	go get -u -d -v ./...
	go mod tidy -v

test:
	go vet ./...
	go test -coverprofile cover.out ./...
	golangci-lint run --fix ./...

add_test:
	git ls-files | grep '.*\.go$$' | grep -v "_test\.go$$" | while read -r f; do gotests -w -parallel -exported $$f; done

gen_certs:
	openssl ecparam -name prime256v1 -genkey > pkg/client/builder/certs/tmp.key
	openssl req -new -sha256 -key pkg/client/builder/certs/tmp.key > pkg/client/builder/certs/tmp.csr
	openssl x509 -req -signkey pkg/client/builder/certs/tmp.key < pkg/client/builder/certs/tmp.csr > pkg/client/builder/certs/tmp.crt

dep:
	go mod tidy
	go mod vendor

unittest:
	cd pkg
	go test -v ./...
build:
	go build ./cmd/boost

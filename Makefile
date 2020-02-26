.PHONY: install

build:
	go build -o nitro .
run: build
	./nitro init
clean:
	multipass delete nitro-dev
	multipass purge
test:
	go test ./...
install:
	go install
release: test
	GOOS=linux GOARCH=amd64 go build -o bin/nitro-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/nitro-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o bin/nitro-windows-amd64

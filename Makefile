.PHONY: install

MACHINE ?= demo-app

build:
	go build -o nitro ./cmd/nitro
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
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/nitro-linux-amd64 ./cmd/nitro
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/nitro-darwin-amd64 ./cmd/nitro
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/nitro-windows-amd64 ./cmd/nitro
demo: build
	composer create-project craftcms/craft ${MACHINE}
	./nitro --machine ${MACHINE} init
	./nitro --machine ${MACHINE} add-host ${MACHINE}.nitro
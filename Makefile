.PHONY: install

build:
	go build -o nitro .
run: build
	./nitro init
clean:
	multipass delete nitro
	multipass purge
test:
	go test ./...
install:
	go install
release: test
	go build -o bin/nitro

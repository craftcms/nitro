.PHONY: install

build:
	go build -o phpdev .
run: build
	./phpdev init
clean:
	multipass delete phpdev
	multipass purge
test:
	go test ./...
install:
	go install

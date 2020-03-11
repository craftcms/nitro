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
demo: build
	composer create-project craftcms/craft ${MACHINE}
	./nitro --machine ${MACHINE} init
	./nitro --machine ${MACHINE} add-host demo ${MACHINE}
releaser:
	goreleaser --snapshot --skip-publish --rm-dist
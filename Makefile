.PHONY: install

MACHINE ?= demo-app

build:
	go build -o nitro ./cmd/next
run: build
	./nitro init
clean:
	multipass delete nitro-dev
	multipass purge
test:
	go test ./...
demo: build
	composer create-project craftcms/craft demo-site
	./nitro --machine ${MACHINE} init
	./nitro --machine ${MACHINE} site --path=demo-site demo
releaser:
	goreleaser --snapshot --skip-publish --rm-dist
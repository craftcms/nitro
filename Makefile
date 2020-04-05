.PHONY: install

MACHINE ?= nitro-global

build:
	go build -o nitro ./cmd/cli
run: build
	./nitro init
clean:
	multipass delete nitro-dev
	multipass purge
test:
	go test ./...
demo: build
	./nitro serve demo-site demo
demo-site:
	composer create-project craftcms/craft demo-site
releaser:
	goreleaser --snapshot --skip-publish --rm-dist
integration-test: build
	./nitro -f nitro.yaml machine create
	if [ -d $directory]; then
		@echo "skipping composer project creation step"
	else
		composer create-project craftcms/craft demo-site
	fi
	./nitro -f nitro.yaml serve ./demo-site demo.test

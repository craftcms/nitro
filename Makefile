.PHONY: install

MACHINE ?= nitro-global

build:
	go build -o nitro ./cmd/cli
test:
	go test ./...
demo-site:
	composer create-project craftcms/craft demo-site
releaser:
	goreleaser --skip-publish --rm-dist --skip-validate
integration-test: build
	./nitro -f nitro.yaml machine create
	composer create-project craftcms/craft demo-site
	sudo ./nitro -f nitro.yaml hosts
remove-integration-test:
	./nitro -f nitro.yaml machine destroy -p
	rm -rf demo-site

test-version: releaser
	chmod u+x ./dist/nitro_darwin_amd64/nitro
	./dist/nitro_darwin_amd64/nitro version


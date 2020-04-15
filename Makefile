.PHONY: install
VERSION ?= 0.8.0
MACHINE ?= nitro-global

build:
	go build -ldflags="-s -w -X github.com/pixelandtonic/nitro/internal/cmd/cmd.Version=${VERSION}" -o nitro ./cmd/cli
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

test-version: build
	./nitro version


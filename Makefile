.PHONY: install
VERSION ?= 1.0.0-beta.4

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro ./cmd/cli
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro.exe ./cmd/cli
local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe
test:
	go test ./...
releaser:
	goreleaser --skip-publish --rm-dist --skip-validate
win-home:
	mkdir "${HOME}"/Nitro

.PHONY: install scripts docker

VERSION ?= 2.0.0-alpha
NITRO_DEFAULT_MACHINE ?= nitro-dev

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro ./cmd/cli
build-api:
	GOOS=linux go build -ldflags="-s -w" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro.exe ./cmd/cli

docker:
	docker build -t craftcms/nitro-proxy:develop .

local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe

dev: scripts api

test:
	go test ./...

vet:
	go vet ./...

releaser:
	goreleaser --skip-publish --rm-dist --skip-validate

win-home:
	mkdir "${HOME}"/Nitro

proto:
	protoc protob/nitro.proto --go_out=plugins=grpc:.

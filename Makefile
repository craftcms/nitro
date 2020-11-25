.PHONY: docker

VERSION ?= 2.0.0-alpha
NITRO_DEFAULT_MACHINE ?= nitro-dev

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/cli
build-api:
	GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/cli

docker:
	docker build -t craftcms/nitro-proxy:${VERSION} .

local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe

dev: api

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

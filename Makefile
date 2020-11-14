.PHONY: install scripts docker

VERSION ?= 1.1.0
NITRO_DEFAULT_MACHINE ?= nitro-dev

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro ./cmd/cli
build-api:
	GOOS=linux go build -ldflags="-s -w" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro.exe ./cmd/cli

build-v2:
	go build -ldflags="-s -w" -o nitro ./cmd/v2

docker:
	docker build -t craftcms/nitro-proxy:develop .

local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe
v2-local: build-v2
	mv nitro /usr/local/bin/nitro

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
	protoc internal/nitrod/nitrod.proto --go_out=plugins=grpc:.
v2-proto:
	protoc pkg/protob/nitro.proto --go_out=plugins=grpc:.

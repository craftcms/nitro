.PHONY: docker

VERSION ?= 2.0.0-alpha

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-api:
	GOOS=linux go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro

docker:
	docker build --build-arg NITRO_VERSION=${VERSION} -t craftcms/nitro-proxy:${VERSION} .

local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe
local-linux: build
	mv nitro ${HOME}/bin/nitro

dev: rm docker init
rm:
	docker container rm -f nitro-v2
init:
	nitro init

test:
	go test ./...
coverage:
	go test ./... -coverprofile profile.out
	go tool cover -html=profile.out
vet:
	go vet ./...

releaser:
	goreleaser --skip-publish --rm-dist --skip-validate

win-home:
	mkdir "${HOME}"/Nitro

proto:
	protoc protob/nitro.proto --go_out=plugins=grpc:.

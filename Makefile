.PHONY: docker docs

VERSION ?= 2.0.0-alpha

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-api:
	GOOS=linux go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro.exe ./cmd/nitro
upx: build
	upx --brute nitro

alpha: alpha-macos alpha-win
alpha-macos: build
	zip -X macos_nitro_v2_alpha.zip nitro
	rm nitro
alpha-win: build-win
	zip -X windows_nitro_v2_alpha.zip nitro
	rm nitro

docker:
	docker build --build-arg NITRO_VERSION=${VERSION} -t craftcms/nitro-proxy:${VERSION} .
docs:
	go run cmd/docs/main.go

local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe
local-linux: build
	sudo mv nitro /usr/local/bin/nitro
local-prod: build upx
	mv nitro /usr/local/bin/nitro

dev: rm docker init
rm:
	docker container rm -f nitro-v2
init:
	nitro init

test:
	go test ./...
coverage:
	go test -v ./... -coverprofile profile.out
	go tool cover -html=profile.out
vet:
	go vet ./...

releaser:
	goreleaser --skip-publish --rm-dist --skip-validate

win-home:
	mkdir "${HOME}"/Nitro

proto:
	protoc protob/nitro.proto --go_out=plugins=grpc:.

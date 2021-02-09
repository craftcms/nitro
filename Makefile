.PHONY: docker docs

VERSION ?= 2.0.0-beta.3

build:
	go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-macos:
	GOOS=darwin go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-macos-arm:
	GOOS=darwin GOARCH=arm64 go1.16rc1 build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-api:
	go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro.exe ./cmd/nitro
build-linux:
	GOOS=linux go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-linux-arm:
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
upx: build
	upx --brute nitro

beta: beta-macos beta-macos-arm beta-win beta-linux beta-linux-arm
beta-macos: build-macos
	tar -cvzf nitro_darwin_x84_64.tar.gz nitro
	rm nitro
beta-macos-arm: build-macos-arm
	tar -cvzf nitro_darwin_arm64.tar.gz nitro
	rm nitro
beta-win: build-win
	zip -X nitro_windows_x86_64.zip nitro.exe
	rm nitro.exe
beta-linux: build-linux
	tar -cvzf nitro_linux_x86_64.tar.gz nitro
	rm nitro
beta-linux-arm: build-linux-arm
	tar -cvzf nitro_linux_arm64.tar.gz nitro
	rm nitro

mod:
	go mod tidy && go mod verify

upx-macos:
	upx --brute nitro
upx-win:
	upx --brute nitro.exe
upx-linux:
	upx --brute nitro

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
	protoc protob/nitrod.proto --go_out=plugins=grpc:.

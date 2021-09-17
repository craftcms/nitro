.PHONY: proxy docs

VERSION ?= 3.0.0

build:
	go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-macos:
	GOOS=darwin go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-api:
	go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro.exe ./cmd/nitro
build-linux:
	GOOS=linux go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro
build-linux-arm:
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w -X 'github.com/craftcms/nitro/command/version.Version=${VERSION}'" -o nitro ./cmd/nitro

mod:
	go mod tidy && go mod verify

proxy:
	docker build --build-arg NITRO_VERSION=${VERSION} -t craftcms/nitro-proxy:${VERSION} .

docs:
	go run cmd/docs/main.go
	
images:
	cd image && $(MAKE) all

local: build
	mv nitro /usr/local/bin/nitro
local-linux: build
	sudo mv nitro /usr/local/bin/nitro

test:
	go test -v ./...
coverage:
	go test -v ./... -coverprofile profile.out
	go tool cover -html=profile.out
vet:
	go vet ./...

releaser:
	goreleaser --skip-publish --rm-dist --skip-validate

proto:
	protoc protob/nitrod.proto --go_out=plugins=grpc:.

.PHONY: install scripts

VERSION ?= 1.0.2
NITRO_DEFAULT_MACHINE ?= nitro-dev

build:
	go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro ./cmd/cli
build-api:
	GOOS=linux go build -ldflags="-s -w" -o nitrod ./cmd/nitrod
build-win:
	GOOS="windows" go build -ldflags="-s -w -X 'github.com/craftcms/nitro/internal/cmd.Version=${VERSION}'" -o nitro.exe ./cmd/cli

local: build
	mv nitro /usr/local/bin/nitro
local-win: build-win
	mv nitro.exe "${HOME}"/Nitro/nitro.exe

dev: scripts api

test:
	go test ./...

releaser:
	goreleaser --skip-publish --rm-dist --skip-validate

win-home:
	mkdir "${HOME}"/Nitro

api: build-api
	multipass transfer nitrod ${NITRO_DEFAULT_MACHINE}:/home/ubuntu/nitrod
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl stop nitrod
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo cp /home/ubuntu/nitrod /usr/sbin/
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo chmod u+x /usr/sbin/nitrod
	multipass transfer nitrod.service ${NITRO_DEFAULT_MACHINE}:/home/ubuntu/nitrod.service
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo cp /home/ubuntu/nitrod.service /etc/systemd/system/
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl daemon-reload
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl start nitrod
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl enable nitrod
setup: build-api
	multipass transfer nitrod ${NITRO_DEFAULT_MACHINE}:/home/ubuntu/nitrod
	multipass transfer nitrod.service ${NITRO_DEFAULT_MACHINE}:/home/ubuntu/nitrod.service
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo cp /home/ubuntu/nitrod.service /etc/systemd/system/
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl daemon-reload
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl start nitrod
	multipass exec ${NITRO_DEFAULT_MACHINE} -- sudo systemctl enable nitrod
proto:
	protoc internal/nitrod/nitrod.proto --go_out=plugins=grpc:.
journalctl:
	multipass exec ${NITRO_DEFAULT_MACHINE} -- journalctl -u nitrod -f
scripts:
	multipass mount scripts ${NITRO_DEFAULT_MACHINE}:/home/ubuntu/scripts

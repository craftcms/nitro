.PHONY: install
VERSION ?= 1.0.0-beta.11

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
api-build:
	GOOS=linux go build -ldflags="-s -w" -o nitrod ./cmd/nitrod
api: api-build
	multipass transfer nitrod nitro-dev:/home/ubuntu/nitrod
	multipass exec nitro-dev -- sudo systemctl stop nitrod
	multipass exec nitro-dev -- sudo cp /home/ubuntu/nitrod /usr/sbin/
	multipass exec nitro-dev -- sudo chmod u+x /usr/sbin/nitrod
	multipass transfer cmd/nitrod/nitrod.service nitro-dev:/home/ubuntu/nitrod.service
	multipass exec nitro-dev -- sudo cp /home/ubuntu/nitrod.service /etc/systemd/system/
	multipass exec nitro-dev -- sudo systemctl daemon-reload
	multipass exec nitro-dev -- sudo systemctl start nitrod
setup: api-build
	multipass transfer nitrod nitro-dev:/home/ubuntu/nitrod
	multipass transfer cmd/nitrod/nitrod.service nitro-dev:/home/ubuntu/nitrod.service
	multipass exec nitro-dev -- sudo cp /home/ubuntu/nitrod.service /etc/systemd/system/
	multipass exec nitro-dev -- sudo systemctl daemon-reload
	multipass exec nitro-dev -- sudo systemctl start nitrod

IMAGE_USER ?= nitro
VERSION ?= 2.1.0
IMAGE ?= craftcms/nitro:${VERSION}

build:
	docker build --build-arg IMAGE_USER=${IMAGE_USER} --tag ${IMAGE} .
local:
	docker build --build-arg IMAGE_USER=${IMAGE_USER} --tag testing:local .
bash: build
	docker run --name testing --rm -it testing:local bash
run: build
	docker run --name testing -p 8080:80 --rm -it testing:local
exec:
	docker exec -it testing bash
exec-root:
	docker exec -it --user root testing bash

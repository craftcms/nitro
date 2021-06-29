IMAGE_USER ?= nitro
VERSION ?= 3.0.0
IMAGE ?= craftcms/nitro:${VERSION}

build:
	docker build --build-arg IMAGE_USER=${IMAGE_USER} --tag ${IMAGE} .
local:
	docker build --build-arg IMAGE_USER=${IMAGE_USER} --tag nitro:local .
bash: build
	docker run --name nitro --rm -it nitro:local bash
run: local
	docker run --name nitro -p 8080:80 --rm -it nitro:local
exec:
	docker exec -it nitro bash
exec-root:
	docker exec -it --user root nitro bash
push: build
	docker build --build-arg IMAGE_USER=${IMAGE_USER} --tag ${IMAGE} .

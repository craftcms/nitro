IMAGE_USER ?= nitro
VERSION ?= 8.0
IMAGE ?= craftcms/nitro:${VERSION}

build:
	docker build --build-arg PHP_VERSION=${VERSION} --build-arg IMAGE_USER=${IMAGE_USER} --tag ${IMAGE} .
all:
	docker build --build-arg PHP_VERSION=8.0 --build-arg IMAGE_USER=${IMAGE_USER} --tag craftcms/nitro:8.0 .
	docker build --build-arg PHP_VERSION=7.4 --build-arg IMAGE_USER=${IMAGE_USER} --tag craftcms/nitro:7.4 .
	docker build --build-arg PHP_VERSION=7.3 --build-arg IMAGE_USER=${IMAGE_USER} --tag craftcms/nitro:7.3 .
	docker build --build-arg PHP_VERSION=7.2 --build-arg IMAGE_USER=${IMAGE_USER} --tag craftcms/nitro:7.2 .
	docker build --build-arg PHP_VERSION=7.1 --build-arg IMAGE_USER=${IMAGE_USER} --tag craftcms/nitro:7.1 .
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

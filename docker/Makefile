ORG ?= craftcms
TAG ?= develop

deploy: build tag push
build:
	docker build -t nitro-proxy:$(TAG) .
tag:
	docker tag nitro-proxy:$(TAG) "$(ORG)/nitro-proxy:$(TAG)"
push:
	docker push craftcms/nitro-proxy:$(TAG)

network:
	docker network create -d bridge nitro-dev
volume:
	docker volume create nitro-dev

run: build
	docker run --rm --net=nitro-dev -p 80:80 -p 443:443 -p 2019:2019 nitro-proxy:local
run-sh: build
	docker run --rm --net=nitro-dev -p 80:80 -p 443:443 -p 2019:2019 nitro-proxy:local

craft:
	composer create-project craftcms/craft craft-install
craft-run:
	docker run -d --name example.nitro --hostname example.nitro --rm --net=nitro-dev --label com.craftcms.nitro.environment=nitro-dev --label com.craftcms.nitro.host=example.nitro -v=/Users/jasonmccallister/go/src/github.com/craftcms/nitro/docker/craft-install:/app/ craftcms/php-fpm:7.4-dev
craft-inspect:
	docker container inspect local-craft
craft-exec:
	docker exec -it local-craft sh
craft-kill:
	docker kill local-craft

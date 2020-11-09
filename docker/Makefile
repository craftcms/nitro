build:
	docker build -t testing-caddy .
network:
	docker network create -d bridge nitro-dev
volume:
	docker volume create nitro-dev
run: build
	docker run --rm --net=nitro-dev -p 80:80 -p 443:443 -p 2019:2019 testing-caddy:latest
run-sh: build
	docker run --rm --net=nitro-dev -p 80:80 -p 443:443 -p 2019:2019 testing-caddy:latest
craft:
	composer create-project craftcms/craft craft-install
craft-run:
	docker run -d --name nitro-dev --hostname nitro-dev --rm --net=nitro-dev --label com.craftcms.nitro.site=nitro-dev -v=/Users/jasonmccallister/go/src/github.com/craftcms/nitro/docker/craft-install:/app/ craftcms/php-fpm:7.4-dev
craft-inspect:
	docker container inspect local-craft
craft-exec:
	docker exec -it local-craft sh
craft-kill:
	docker kill local-craft

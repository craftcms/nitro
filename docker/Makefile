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

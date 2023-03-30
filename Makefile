test-common:
	go test -v ./common/...

test-api:
	go test -v ./api/...

test-distributor:
	go test -v ./distributor/...

test-indexer:
	go test -v ./indexer/...

build-common:
	go build -o bin/common ./common/

build-api:
	go build -o bin/api ./api/

build-distributor:
	go build -o bin/distributor ./distributor/

build-indexer:
	go build -o bin/indexer ./indexer/

run-api:
	ENV=dev go run api/main.go

run-distributor:
	ENV=dev go run distributor/main.go

run-indexer:
	ENV=dev go run indexer/main.go

docker-build-api:
	docker build -f .docker/DockerfileApi  -t api .
	
docker-build-distributor:
	docker build -f .docker/DockerfileDistributor  -t distributor .

docker-build-indexer:
	docker build -f .docker/DockerfileIndexer  -t indexer .

docker-run-api:
	docker rm -f api && docker run -d -p 8080:8080 --env-file .env/.env.api.docker --name api api

docker-run-distributor:
	docker rm -f distributor && docker run -d --env-file .env/.env.distributor.docker --name distributor distributor

docker-run-indexer:
	docker rm -f indexer && docker run -d --env-file .env/.env.indexer.docker --name indexer indexer

docker-run-postgres:
	docker rm -f postgres && docker run -e POSTGRES_PASSWORD=mysecretpassword -d --name postgres postgres

docker-run-redis:
	docker rm -f redis && docker run -d -p 6379:6379 --name redis redis

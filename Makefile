test-common:
	go test -v ./common/...

test-domain:
	go test -v ./domain/...

test-api:
	go test -v ./cmd/api/...

test-distributor:
	go test -v ./cmd/distributor/...

test-indexer:
	go test -v ./cmd/indexer/...

test-migrator:
	go test -v ./cmd/migrator/...

test-ethereum:
	go test -v ./gateways/ethereum/...

test-redis:
	go test -v ./gateways/redis/...

test-images:
	go test -v ./gateways/images/...

test-postgres:
	go test -v ./gateways/postgres/...

run-api:
	ENV=dev APP_NAME=api go run cmd/api/*.go

run-distributor:
	ENV=dev APP_NAME=distributor go run cmd/distributor/*.go

run-indexer:
	ENV=dev APP_NAME=indexer go run cmd/indexer/*.go

run-migrator:
	ENV=dev APP_NAME=migrator go run cmd/migrator/*.go

postgres-add:
	goose -dir gateways/postgres/migrations create $(name).sql

postgres-bindings:
	sqlc generate --experimental -f "gateways/postgres/sqlc.yml"

docker-build-api:
	docker build -f .docker/DockerfileApi  -t api .
	
docker-build-distributor:
	docker build -f .docker/DockerfileDistributor  -t distributor .

docker-build-indexer:
	docker build -f .docker/DockerfileIndexer  -t indexer .

docker-build-migrator:
	docker build -f .docker/DockerfileMigrator  -t migrator .

docker-run-api:
	docker rm -f api && docker run -d -p 8080:8080 --env-file .env/.env.api.docker --name api api

docker-run-distributor:
	docker rm -f distributor && docker run -d --env-file .env/.env.distributor.docker --name distributor distributor

docker-run-indexer:
	docker rm -f indexer && docker run -d --env-file .env/.env.indexer.docker --name indexer indexer

docker-run-migrator:
	docker rm -f migrator && docker run -d --env-file .env/.env.migrator.docker --name migrator migrator

docker-run-postgres:
	docker rm -f postgres && docker run -d  -p 5432:5432 --env-file .env/.env.postgres.docker --name postgres postgres

docker-run-redis:
	docker rm -f redis && docker run -d -p 6379:6379 --name redis redis

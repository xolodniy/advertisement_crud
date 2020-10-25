#test1
GO111MODULE=on

# build application context: documentation && mocks for integration tests
.PHONY: build
build:
	docker-compose build

.PHONY: up
up:
	docker-compose up -d db
	docker-compose up server

.PHONY: drop_db
drop_db:
	docker-compose rm -sf db

# console access to database
.PHONY: mysql
mysql:
	docker-compose up -d db
	while true ; do sleep 0.1; docker-compose exec db pg_isready && break ; done
	docker-compose exec db mysql -u advertisement_crud -pQEU2zYILBoMAH26Q advertisement_crud

.PHONY: test
test_unit:
	docker-compose run --rm server go test ./...
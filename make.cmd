@ECHO OFF

IF "%1"=="up" (
    docker-compose up server

) ELSE IF "%1"=="build" (
    docker-compose build

) ELSE IF "%1"=="drop_db" (
	docker-compose rm -sf db

) ELSE IF "%1"=="mysql" (
	docker-compose up -d db
	timeout 1 > NUL
	docker-compose exec db mysql -u advertisement_crud -pQEU2zYILBoMAH26Q advertisement_crud

) ELSE IF "%1"=="test" (
    docker-compose run --rm server go test ./...

) ELSE (
	echo unexpected command
)
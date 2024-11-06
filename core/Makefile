sayhello:
	@echo "Hello World"

dockerize:
	docker run --name rating-db-dev -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=developer -e POSTGRES_DB=rating -p 5432:5432 -d postgres:16.3-alpine3.20

undockerize:
	docker rm -f rating-db-dev

migrateup:
	migrate -path database/migrations -database "postgresql://developer:secret@localhost:5432/rating?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migrations -database "postgresql://developer:secret@localhost:5432/rating?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

gomock:
	mockgen -package mock -destination internal/repository/mock/querier.go github.com/zepyrshut/rating-orama/internal/repository ExtendedQuerier

run:
	go run ./cmd/.

recreate:
	make undockerize
	make dockerize
	sleep 2
	make migrateup

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./tmp/arena ./cmd/.

pack-docker:
	make test
	make build-linux
	docker build -t rating:${version} -t rating:latest .

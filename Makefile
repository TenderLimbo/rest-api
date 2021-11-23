build:
	docker-compose build restapi

run:
	docker-compose up restapi

test:
	go test -v ./...

lint:
	golangci-lint run



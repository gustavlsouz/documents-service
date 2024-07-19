install:
	go mod download

start: 
	docker compose --file ./deployments/docker-compose.yaml up -d  && cd cmd && go run ./main.go

init:
	docker compose --file ./deployments/docker-compose.yaml up -d

run-test:
	go test -v ./...
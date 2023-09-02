test:
	go test -v ./... -coverprofile=c.out

cover:
	go tool cover -html=c.out

dev:
	air

build:
	go build -o bin/main main.go

run:
	make build && ./bin/main
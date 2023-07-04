FROM golang:alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go get ./...

CMD ["go", "run", "cmd/main.go"]
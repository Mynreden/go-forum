FROM golang:latest

workdir /app

copy go.mod go.sum ./

run go mod download

copy . . 

run go build -o main ./cmd/

RUN apt-get update && apt-get install -y sqlite3

expose 8000

CMD ["./main"]
FROM golang:alpine

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum crud-server.go /app/
RUN go mod download
RUN go build -o main .
ENTRYPOINT ["./main"]

FROM golang:1.18-buster AS gobuilder

ENV CGO_ENABLED 0

COPY . /app
WORKDIR /app

RUN go build main.go
CMD ["./main"]
EXPOSE 8080
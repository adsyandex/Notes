# Dockerfile
FROM golang:1.20-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o notes_app

CMD ["./notes_app"]

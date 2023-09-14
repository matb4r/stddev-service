FROM golang:1.21.0
LABEL authors="mateusz.bartkowiak"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /stddev-service

EXPOSE 80

CMD ["/stddev-service"]
FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run github.com/99designs/gqlgen generate
RUN CGO_ENABLED=0 GOOS=linux go build -o ./lapis

EXPOSE 8080

CMD ["./lapis"]

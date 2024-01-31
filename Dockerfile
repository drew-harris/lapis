FROM node:18 as tailwind

WORKDIR /app
COPY . .
RUN npm install
RUN npm run twgen

FROM golang:1.21

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go run github.com/99designs/gqlgen generate

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go mod tidy
RUN templ generate

COPY --from=tailwind app/dist/out.css ./dist/out.css


RUN CGO_ENABLED=0 GOOS=linux go build -o ./lapis

EXPOSE 8080

CMD ["./lapis"]

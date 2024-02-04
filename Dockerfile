FROM node:18 as javascript

WORKDIR /app
COPY . .
RUN npm install -g pnpm
RUN pnpm i
RUN pnpm run tw:gen
RUN pnpm run build

FROM golang:1.21

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

RUN go run github.com/99designs/gqlgen generate

RUN go mod tidy

COPY --from=javascript app/dist/out.css ./dist/out.css
COPY --from=javascript app/dist/client/index.js ./dist/client/index.js


RUN CGO_ENABLED=0 GOOS=linux go build -o ./lapis

EXPOSE 8080

CMD ["./lapis"]

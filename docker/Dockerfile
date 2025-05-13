FROM golang:1.24.1-alpine as build

WORKDIR /app

COPY .. .

RUN go mod download  \
    && go mod tidy

RUN go build -o main cmd/service.go

FROM debian as app

WORKDIR /app

COPY --from=build /app/main main
COPY --from=build /app/migrations migrations
COPY --from=build /app/config config

ENV APP_ENV=docker

EXPOSE 50052

CMD ["./main"]
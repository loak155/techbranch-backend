# Build stage
FROM golang:1.21.5-alpine3.18 AS development
WORKDIR /app
COPY . .
RUN go build -o main ./cmd

# Run stage
FROM alpine:3.18 AS production
WORKDIR /app
COPY --from=development /app/main .
COPY .env .
COPY migrations ./migrations

EXPOSE 8080 9090
CMD [ "/app/main" ]
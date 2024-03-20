# Build stage
FROM golang:1.21.5-alpine3.18 AS development
WORKDIR /app
COPY . .
RUN go install github.com/cosmtrek/air@v1.49.0
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2
RUN go build -o main ./cmd/article

# Run stage
FROM alpine:3.18 AS production
WORKDIR /app
COPY configs/ configs/
COPY --from=development /app/main .
CMD [ "/app/main" ]
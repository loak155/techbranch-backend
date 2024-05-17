# Build stage
FROM golang:1.21.5-alpine3.18 AS development
WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -o main ./services/bookmark/cmd

# Run stage
FROM alpine:3.18 AS production
WORKDIR /app
COPY --from=development /app/main .
CMD [ "/app/main" ]
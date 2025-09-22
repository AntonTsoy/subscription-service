FROM golang:1.23.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init -g cmd/app/main.go -o ./docs
RUN GOOS=linux go build -o server ./cmd/app


FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
COPY ./migrations ./migrations
COPY ./config ./config
CMD ["./server"]

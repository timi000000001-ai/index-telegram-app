FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot-service ./main.go

FROM gcr.io/distroless/base

COPY --from=builder /app/bot-service /bot-service

EXPOSE 8080

CMD ["/bot-service"]
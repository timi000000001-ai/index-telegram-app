FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o collection-service ./main.go

FROM gcr.io/distroless/base

COPY --from=builder /app/collection-service /collection-service

EXPOSE 8080

CMD ["/collection-service"]
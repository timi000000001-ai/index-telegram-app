FROM golang:1.25 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o management-service ./main.go

FROM gcr.io/distroless/base

COPY --from=builder /app/management-service /management-service

EXPOSE 8080

CMD ["/management-service"]
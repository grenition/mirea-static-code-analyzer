FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cpp_analyzer_service ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates cppcheck

WORKDIR /app

COPY --from=builder /app/cpp_analyzer_service .

EXPOSE 8085

CMD ["./cpp_analyzer_service"]


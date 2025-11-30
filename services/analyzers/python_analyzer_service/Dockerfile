FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/python_analyzer_service ./main.go

FROM python:3.11-alpine

RUN apk --no-cache add ca-certificates && \
    pip install flake8

WORKDIR /app

COPY --from=builder /app/python_analyzer_service .

EXPOSE 8082

CMD ["./python_analyzer_service"]


FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/javascript_analyzer_service ./main.go

FROM node:18-alpine

RUN apk --no-cache add ca-certificates && \
    npm install -g eslint

WORKDIR /app

COPY --from=builder /app/javascript_analyzer_service .

EXPOSE 8083

CMD ["./javascript_analyzer_service"]


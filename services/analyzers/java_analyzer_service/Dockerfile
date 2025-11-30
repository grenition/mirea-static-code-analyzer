FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/java_analyzer_service ./main.go

FROM eclipse-temurin:17-jdk

RUN apt-get update && apt-get install -y wget && \
    wget -O /tmp/checkstyle.jar https://github.com/checkstyle/checkstyle/releases/download/checkstyle-10.12.5/checkstyle-10.12.5-all.jar && \
    echo '#!/bin/bash' > /usr/local/bin/checkstyle && \
    echo 'java -jar /tmp/checkstyle.jar "$@"' >> /usr/local/bin/checkstyle && \
    chmod +x /usr/local/bin/checkstyle && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/java_analyzer_service .

EXPOSE 8084

CMD ["./java_analyzer_service"]


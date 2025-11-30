FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/csharp_analyzer_service ./main.go

FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/csharp_analyzer_service .

EXPOSE 8086

CMD ["./csharp_analyzer_service"]


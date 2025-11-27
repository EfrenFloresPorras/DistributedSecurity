# ---------- Etapa 1: Build ----------
FROM golang:1.24 AS builder

WORKDIR /app

# Copiar dependencias desde la raíz
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el proyecto
COPY . .

# Compilar ThreatLog Service
RUN go build -o /app/bin/threatlog-service ./threatlog-service/cmd/main

# ---------- Etapa 2: Imagen final ----------
FROM debian:bookworm-slim

WORKDIR /app

# Copiar binario desde etapa builder
COPY --from=builder /app/bin/threatlog-service /app/threatlog-service

EXPOSE 8081

CMD ["/app/threatlog-service"]
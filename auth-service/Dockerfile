# ---------- Etapa 1: Build ----------
FROM golang:1.23 AS builder

WORKDIR /app

# Copiar dependencias desde la raíz
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el proyecto
COPY . .

# Compilar Auth Service
RUN go build -o /app/bin/auth-service ./auth-service/cmd/main

# ---------- Etapa 2: Imagen final ----------
FROM debian:bookworm-slim

WORKDIR /app

# Copiar binario desde etapa builder
COPY --from=builder /app/bin/auth-service /app/auth-service

EXPOSE 8080

CMD ["/app/auth-service"]
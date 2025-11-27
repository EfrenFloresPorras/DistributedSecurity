# ---------- Build ----------
FROM golang:1.23 AS builder
WORKDIR /app

# Usa el go.mod/go.sum de la RAÍZ del monorepo como contexto
COPY go.mod go.sum ./
RUN go mod download

# Copia TODO el repo (necesario si importas paquetes internos)
COPY . .

# Compila el binario del GoServer
# Si tu main está en GoServer/main.go, esta ruta es correcta:
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./GoServer

# ---------- Runtime ----------
# Distroless = imagen mínima y segura
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/bin/server /app/server

EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/server"]

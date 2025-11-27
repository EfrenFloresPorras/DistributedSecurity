FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# 👇 cambia la ruta del build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/policy ./policy-service/cmd/main/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/bin/policy /app/policy
EXPOSE 8082
USER 65532:65532
ENTRYPOINT ["/app/policy"]

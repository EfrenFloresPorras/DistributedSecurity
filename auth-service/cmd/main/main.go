package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	consulreg "DistributedSecurity/pkg/discovery/consul"
	"DistributedSecurity/pkg/registry"

	grpcserver "DistributedSecurity/auth-service/grpc" // tu nuevo servidor gRPC
)

func main() {
	fmt.Println("🚀 Auth Service (gRPC) starting on port 8080...")

	// Configuración de Consul
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr == "" {
		consulAddr = "http://consul:8500"
	}
	reg, err := consulreg.NewRegistry(consulAddr)
	if err != nil {
		log.Fatalf("❌ Error conectando a Consul: %v", err)
	}

	// Registrar instancia en Consul
	instanceID := registry.GenerateInstanceID("auth-service")
	ctx := context.Background()
	if err := reg.Register(ctx, instanceID, "auth-service", "auth-service:8080"); err != nil {
		log.Fatalf("❌ No se pudo registrar en Consul: %v", err)
	}
	log.Printf("📌 Auth Service registrado en Consul con ID %s", instanceID)

	// Heartbeat (TTL)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			_ = reg.ReportHealthyState(instanceID, "auth-service")
		}
	}()

	// Iniciar servidor gRPC
	grpcserver.StartGRPCServer()
}

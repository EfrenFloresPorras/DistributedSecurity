package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	consulreg "DistributedSecurity/pkg/discovery/consul"
	"DistributedSecurity/pkg/registry"

	grpcserver "DistributedSecurity/policy-service/grpc"
)

func main() {
	fmt.Println("🚀 Policy Service (gRPC) starting on port 8082...")

	// ----------------------------
	// 1. Conexión a Consul
	// ----------------------------
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr == "" {
		consulAddr = "http://consul:8500"
	}
	reg, err := consulreg.NewRegistry(consulAddr)
	if err != nil {
		log.Fatalf("❌ Error conectando a Consul: %v", err)
	}

	// ----------------------------
	// 2. Registro del servicio
	// ----------------------------
	instanceID := registry.GenerateInstanceID("policy-service")
	ctx := context.Background()
	if err := reg.Register(ctx, instanceID, "policy-service", "policy-service:8082"); err != nil {
		log.Fatalf("❌ No se pudo registrar en Consul: %v", err)
	}
	log.Printf("📌 Policy Service registrado en Consul con ID %s", instanceID)

	// ----------------------------
	// 3. Heartbeat
	// ----------------------------
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			_ = reg.ReportHealthyState(instanceID, "policy-service")
		}
	}()

	// ----------------------------
	// 4. Iniciar servidor gRPC
	// ----------------------------
	grpcserver.StartGRPCServer()
}

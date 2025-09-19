package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"DistributedSecurity/auth-service/internal/controller/auth"
	threatlog "DistributedSecurity/auth-service/internal/gateway/threatlog"
	authhttp "DistributedSecurity/auth-service/internal/handler/http"
	"DistributedSecurity/auth-service/internal/repository/memory"
	"DistributedSecurity/auth-service/pkg/model"

	consulreg "DistributedSecurity/pkg/discovery/consul"
	"DistributedSecurity/pkg/registry"
)

func main() {
	fmt.Println("🚀 Auth Service starting on port 8080...")

	// Repositorio en memoria
	repo := memory.New()
	repo.Put(&model.User{Username: "admin", PasswordHash: "1234"})

	// Cliente de Consul
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr == "" {
		consulAddr = "http://consul:8500"
	}
	reg, err := consulreg.NewRegistry(consulAddr)
	if err != nil {
		log.Fatalf("❌ Error conectando a Consul: %v", err)
	}

	// Cliente ThreatLog (usando Consul)
	threatLogger := threatlog.NewClient(reg)

	ctrl := auth.New(repo, threatLogger)
	h := authhttp.New(ctrl)

	http.HandleFunc("/healthz", h.Healthz)
	http.HandleFunc("/login", h.Login)

	// Iniciar servidor HTTP
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Error Auth Service: %v", err)
		}
	}()
	log.Println("✅ Auth Service corriendo en :8080")

	// Registrar en Consul
	instanceID := registry.GenerateInstanceID("auth-service")
	ctx := context.Background()
	if err := reg.Register(ctx, instanceID, "auth-service", "auth-service:8080"); err != nil {
		log.Fatalf("❌ No se pudo registrar en Consul: %v", err)
	}
	log.Printf("📌 Auth Service registrado en Consul con ID %s", instanceID)

	// Heartbeat
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		_ = reg.ReportHealthyState(instanceID, "auth-service")
	}
}

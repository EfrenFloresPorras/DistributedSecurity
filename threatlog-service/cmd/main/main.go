package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"DistributedSecurity/threatlog-service/internal/controller/threat"
	threathttp "DistributedSecurity/threatlog-service/internal/handler/http"
	"DistributedSecurity/threatlog-service/internal/repository/memory"

	consulreg "DistributedSecurity/pkg/discovery/consul"
	"DistributedSecurity/pkg/registry"
)

func main() {
	fmt.Println("Threat Logging Service starting on port 8081...")

	// ----------------------------
	// 1. Repositorio + controlador
	// ----------------------------
	repo := memory.New()
	ctrl := threat.New(repo)
	h := threathttp.New(ctrl)

	// ----------------------------
	// 2. Endpoints HTTP
	// ----------------------------
	http.HandleFunc("/healthz", h.Healthz)
	http.HandleFunc("/log", h.Log)
	http.HandleFunc("/events", h.Events)

	// ----------------------------
	// 3. Servidor HTTP
	// ----------------------------
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("Error starting ThreatLog Service: %v", err)
		}
	}()
	log.Println("ThreatLog Service running at :8081")

	// ----------------------------
	// 4. Registro en Consul
	// ----------------------------
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr == "" {
		consulAddr = "http://localhost:8500"
	}
	reg, err := consulreg.NewRegistry(consulAddr)
	if err != nil {
		log.Fatalf("Error connecting to Consul: %v", err)
	}

	instanceID := registry.GenerateInstanceID("threatlog-service")
	ctx := context.Background()
	if err := reg.Register(ctx, instanceID, "threatlog-service", "threatlog-service:8081"); err != nil {
		log.Fatalf("Failed to register in Consul: %v", err)
	}
	log.Printf("ThreatLog Service registered in Consul with ID %s", instanceID)

	// ----------------------------
	// 5. Heartbeat
	// ----------------------------
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		_ = reg.ReportHealthyState(instanceID, "threatlog-service")
	}
}

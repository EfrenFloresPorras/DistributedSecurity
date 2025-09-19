package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"DistributedSecurity/pkg/registry"

	api "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *api.Client
}

// NewRegistry crea un nuevo cliente Consul
func NewRegistry(addr string) (*Registry, error) {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

// Register registra una nueva instancia en Consul
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort debe tener formato host:puerto, ej: localhost:8080")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: parts[0],
		Port:    port,
		Check: &api.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
}

// Deregister elimina una instancia de Consul
func (r *Registry) Deregister(ctx context.Context, instanceID string, _ string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddress devuelve direcciones activas para un servicio
func (r *Registry) ServiceAddress(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, registry.ErrNotFound
	}

	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

// ReportHealthyState avisa a Consul que la instancia está sana
func (r *Registry) ReportHealthyState(instanceID string, _ string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}

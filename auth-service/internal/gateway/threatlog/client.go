package threatlog

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	consulreg "DistributedSecurity/pkg/discovery/consul"
)

// Client se conecta a Consul y luego a ThreatLog
type Client struct {
	registry *consulreg.Registry
}

func NewClient(reg *consulreg.Registry) *Client {
	return &Client{registry: reg}
}

// LogFailedLogin busca threatlog-service en Consul y hace POST /log
func (c *Client) LogFailedLogin(ctx context.Context, username string) error {
	// Buscar instancias en Consul
	addresses, err := c.registry.ServiceAddress(ctx, "threatlog-service")
	if err != nil {
		return fmt.Errorf("no se encontró threatlog-service en Consul: %w", err)
	}
	if len(addresses) == 0 {
		return fmt.Errorf("no hay instancias de threatlog-service disponibles")
	}

	// Tomar la primera dirección
	url := fmt.Sprintf("http://%s/log", addresses[0])
	body := []byte(fmt.Sprintf(`{"type":"failed_login","username":"%s"}`, username))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error al llamar a ThreatLog: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ThreatLog devolvió %d", resp.StatusCode)
	}
	return nil
}

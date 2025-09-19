package registry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry define las operaciones de un servicio de discovery
type Registry interface {
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	ServiceAddress(ctx context.Context, serviceName string) ([]string, error)
	ReportHealthyState(instanceID string, serviceName string) error
}

// Error cuando no hay instancias
var ErrNotFound = errors.New("no service addresses found")

// GenerateInstanceID crea un ID único para cada instancia
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

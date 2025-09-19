package memory

import (
	"DistributedSecurity/threatlog-service/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	events []model.ThreatEvent
}

func New() *Repository {
	return &Repository{events: []model.ThreatEvent{}}
}

func (r *Repository) Add(event model.ThreatEvent) {
	r.Lock()
	defer r.Unlock()
	r.events = append(r.events, event)
}

func (r *Repository) List() []model.ThreatEvent {
	r.RLock()
	defer r.RUnlock()
	return append([]model.ThreatEvent{}, r.events...) // copia
}

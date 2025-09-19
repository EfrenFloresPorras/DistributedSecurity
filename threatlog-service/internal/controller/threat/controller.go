package threat

import (
	"DistributedSecurity/threatlog-service/pkg/model"
	"time"
)

type threatRepository interface {
	Add(event model.ThreatEvent)
	List() []model.ThreatEvent
}

type Controller struct {
	repo threatRepository
}

func New(repo threatRepository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) LogEvent(eventType, ip string) model.ThreatEvent {
	e := model.ThreatEvent{
		ID:        time.Now().Format("20060102150405"), // timestamp como ID
		Type:      eventType,
		IP:        ip,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	c.repo.Add(e)
	return e
}

func (c *Controller) ListEvents() []model.ThreatEvent {
	return c.repo.List()
}

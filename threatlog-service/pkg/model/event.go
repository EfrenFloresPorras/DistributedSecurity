package model

type ThreatEvent struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
}

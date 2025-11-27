package policy

import "log"

// Simple controller para extender lógica futura
func ValidateUser(username string) bool {
	log.Printf("[Policy Controller] Validating user=%s", username)
	return username != "blocked"
}

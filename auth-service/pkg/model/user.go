package model

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type Token struct {
	Value      string `json:"value"`
	Expiration string `json:"expiration"`
}

// Payload de login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

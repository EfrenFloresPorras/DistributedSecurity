package auth

import (
	"DistributedSecurity/auth-service/pkg/model"
	"context"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type authRepository interface {
	Get(username string) (*model.User, bool)
	Put(u *model.User)
}

type threatLogger interface {
	LogFailedLogin(ctx context.Context, username string) error
}

type Controller struct {
	repo   authRepository
	logger threatLogger
}

func New(repo authRepository, logger threatLogger) *Controller {
	return &Controller{repo: repo, logger: logger}
}

func (c *Controller) Login(req model.LoginRequest, ctx context.Context) (*model.LoginResponse, error) {
	user, ok := c.repo.Get(req.Username)
	if !ok || user.PasswordHash != req.Password {
		if c.logger != nil {
			_ = c.logger.LogFailedLogin(ctx, req.Username) // registrar en ThreatLog vía Consul
		}
		return nil, ErrInvalidCredentials
	}

	return &model.LoginResponse{Token: "fake-token-123"}, nil
}

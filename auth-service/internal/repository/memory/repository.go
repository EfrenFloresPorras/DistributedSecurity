package memory

import (
	"DistributedSecurity/auth-service/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	users map[string]*model.User
}

func New() *Repository {
	return &Repository{
		users: make(map[string]*model.User),
	}
}

func (r *Repository) Get(username string) (*model.User, bool) {
	r.RLock()
	defer r.RUnlock()
	u, ok := r.users[username]
	return u, ok
}

func (r *Repository) Put(u *model.User) {
	r.Lock()
	defer r.Unlock()
	r.users[u.Username] = u
}

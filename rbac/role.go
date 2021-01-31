package rbac

import "sync"

type Role interface {
	ID() string
	Permit(permission Permission) bool
}

type Roles map[string]Role

type DefaultRole struct {
	sync.RWMutex
	id          string
	permissions Permissions
}

func NewDefaultRole(id string) *DefaultRole {
	return &DefaultRole{id: id, permissions: make(Permissions)}
}

func (role *DefaultRole) ID() string {
	return role.id
}

func (role *DefaultRole) Permit(permission Permission) (permit bool) {
	if permission == nil {
		return
	}
	role.RLock()
	defer role.RUnlock()

	for _, p := range role.permissions {
		if p.Match(permission) {
			permit = true
			break
		}
	}
	return
}

func (role *DefaultRole) Assign(p Permission) error {
	role.Lock()
	role.permissions[p.ID()] = p
	role.Unlock()
	return nil
}

func (role *DefaultRole) Remove(p Permission) error {
	role.Lock()
	delete(role.permissions, p.ID())
	role.Unlock()
	return nil
}

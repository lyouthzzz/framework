package rbac

import (
	"errors"
	"sync"
)

var (
	ErrRoleAlreadyExist = errors.New("role has already existed")
	ErrRoleNotFound     = errors.New("role not found")
)

type RBAC struct {
	mu    sync.RWMutex
	roles Roles
}

func New() *RBAC {
	return &RBAC{roles: make(Roles)}
}

func (rbac *RBAC) GetRole(id string) (Role, error) {
	rbac.mu.Lock()
	defer rbac.mu.Unlock()

	role, ok := rbac.roles[id]
	if !ok {
		return role, ErrRoleNotFound
	}
	return role, nil
}

func (rbac *RBAC) AddRole(role Role) (err error) {
	rbac.mu.Lock()
	defer rbac.mu.Unlock()

	if _, ok := rbac.roles[role.ID()]; !ok {
		rbac.roles[role.ID()] = role
	} else {
		err = ErrRoleAlreadyExist
	}
	return
}

func (rbac *RBAC) RemoveRole(id string) (err error) {
	rbac.mu.Lock()
	defer rbac.mu.Unlock()

	if _, ok := rbac.roles[id]; ok {
		delete(rbac.roles, id)
	} else {
		err = ErrRoleNotFound
	}
	return
}

func (rbac *RBAC) IsGranted(id string, permission Permission) (bool, error) {
	role, err := rbac.GetRole(id)
	if err != nil {
		return false, err
	}
	return role.Permit(permission), nil
}

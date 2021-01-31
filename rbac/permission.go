package rbac

type Permission interface {
	ID() string
	Match(permission Permission) bool
}

type Permissions map[string]Permission

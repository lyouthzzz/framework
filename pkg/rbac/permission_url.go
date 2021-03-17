package rbac

import "strings"

const (
	MethodAll = "*"
)

type UrlPermission struct {
	Id     string
	Path   string
	Method string
}

func NewUrlPermission(Id string, path string, method string) *UrlPermission {
	return &UrlPermission{Id: Id, Path: path, Method: method}
}

func (p UrlPermission) ID() string {
	return p.Id
}

func (p *UrlPermission) Match(permission Permission) (match bool) {
	up, ok := permission.(*UrlPermission)
	if !ok {
		return
	}
	if !strings.HasPrefix(up.Path, p.Path) {
		return
	}
	if p.Method == MethodAll {
		match = true
		return
	}
	for _, method := range strings.Split(p.Method, ",") {
		if method == up.Method {
			match = true
			break
		}
	}
	return
}

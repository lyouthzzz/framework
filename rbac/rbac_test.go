package rbac

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRBAC(t *testing.T) {
	var err error

	rbac := New()

	err = rbac.AddRole(NewDefaultRole("a"))
	require.NoError(t, err)

	err = rbac.AddRole(NewDefaultRole("b"))
	require.NoError(t, err)

	roleC := NewDefaultRole("c")
	_ = rbac.AddRole(roleC)

	err = roleC.Assign(NewUrlPermission("p1", "/public", "*"))
	require.NoError(t, err)

	urlPermission := &UrlPermission{Path: "/public", Method: "POST"}

	permit := roleC.Permit(urlPermission)
	require.Equal(t, permit, true)

	granted, err := rbac.IsGranted("a", urlPermission)
	require.NoError(t, err)
	require.Equal(t, granted, false)

	granted, err = rbac.IsGranted("c", urlPermission)
	require.NoError(t, err)
	require.Equal(t, granted, true)
}

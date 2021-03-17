package gormx

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSoftDeleted(t *testing.T) {
	type user struct {
		Name      string
		DeletedAt DeletedAt
	}
	db, err := Connect(MysqlDns("root", "liuyanggg.123", "localhost", 3306, "go_admin", ""), true)
	require.NoError(t, err)

	err = db.AutoMigrate(&user{})
	require.NoError(t, err)

	err = db.WithContext(context.Background()).Create(&user{Name: "name"}).Error
	require.NoError(t, err)

	err = db.WithContext(context.Background()).Where("name = ?", "name").Delete(&user{}).Error
	require.NoError(t, err)
}

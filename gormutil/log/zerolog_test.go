package log

import (
	"context"
	"github.com/lyouthzzz/framework/gormutil"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"testing"
)

type User struct {
	Name  string
	Email string
}

func TestZeroLog(t *testing.T) {
	db, err := gormutil.Connect(gormutil.MysqlDns("root", "root", "localhost", 3306, "framework", ""), true)
	require.NoError(t, err)

	db.Logger = New()

	logger := log.Logger.With().Str("request_id", "test_id").Logger()
	ctx := logger.WithContext(context.Background())

	tx := db.WithContext(ctx).Find(&User{Name: "11"})
	require.NoError(t, tx.Error)
}

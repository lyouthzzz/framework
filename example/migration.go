package main

import (
	"github.com/lyouthzzz/framework/gormutil"
	"github.com/lyouthzzz/framework/gormutil/log"
	"github.com/lyouthzzz/framework/gormutil/migration"
	"gorm.io/gorm"
)

func main() {
	var migrations = []migration.Migration{
		migration.NewMigration("add user", func(db *gorm.DB) error {
			type User struct {
				Name  string
				Email string
			}
			return db.Migrator().CreateTable(&User{})
		}),
	}
	db, err := gormx.Connect(gormx.MysqlDns("root", "liuyanggg.123", "localhost", 3306, "go_admin", ""), true)
	if err != nil {
		panic(err)
	}

	db.Logger = log.New()
	err = migration.Migrate(db, migrations...)
	if err != nil {
		panic(err)
	}

}

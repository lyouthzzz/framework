package main

import (
	"github.com/lyouthzzz/framework/gormutil"
	"github.com/lyouthzzz/framework/gormutil/migration"
	"gorm.io/gorm"
)

func main() {
	var migrations = []migration.Migration{
		migration.NewMigration("11", func(db *gorm.DB) error {
			if db.Migrator().HasTable("demo") {
				return db.Migrator().DropTable("demo")
			}
			return nil
		}),
	}
	db, err := gormutil.Connect(gormutil.MysqlDns("root", "root", "localhost", 3306, "framework", ""), true)
	if err != nil {
		panic(err)
	}
	err = migration.Migrate(db, migrations...)
	if err != nil {
		panic(err)
	}
}

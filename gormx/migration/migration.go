package migration

import "gorm.io/gorm"

type model struct {
	ID string `gorm:"column:id;type:char(60);primaryKey"`
}

func (model) TableName() string {
	return "migration"
}

type Migration interface {
	ID() string
	Migrate(db *gorm.DB) error
}

type migration struct {
	id      string
	migrate func(db *gorm.DB) error
}

func (m *migration) ID() string {
	return m.id
}

func (m *migration) Migrate(db *gorm.DB) error {
	return m.migrate(db)
}

func NewMigration(id string, migrate func(db *gorm.DB) error) Migration {
	return &migration{id: id, migrate: migrate}
}

func Migrate(db *gorm.DB, migrations ...Migration) error {
	var m model
	if !db.Migrator().HasTable(&m) {
		if err := db.Migrator().CreateTable(&m); err != nil {
			return err
		}
	}
	for _, migration := range migrations {
		mm := model{ID: migration.ID()}

		tx := db.Find(&mm)
		if err := tx.Error; err != nil {
			return err
		}
		if tx.RowsAffected == 0 {
			if err := migration.Migrate(db); err != nil {
				return err
			}
			if err := db.Create(&mm).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

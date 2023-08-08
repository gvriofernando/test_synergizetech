package dbmigration

import (
	"github.com/gvriofernando/test_synergizetech/app/store"
	"gorm.io/gorm"
)

type Config struct {
	Db *gorm.DB
}

func DBMigrate(cfg Config) error {
	err := cfg.Db.AutoMigrate(&store.Player{})
	if err != nil {
		return err
	}
	return nil
}

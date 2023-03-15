package db

import (
	"fmt"
	"github.com/potatowhite/books/file-service/cmd/config"
	"github.com/potatowhite/books/file-service/pkg/repository/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {

	dsn := "host=%s port=%v user=%s dbname=%s password=%s sslmode=disable"
	dsn = fmt.Sprintf(dsn, cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Dbname, cfg.Database.Password)

	log.Printf("dsn: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// auto migrate
	autoMigration(err, db)

	return db, nil
}

func CloseDB(database *gorm.DB) {
	func(con *gorm.DB) {
		sqlDB, err := con.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}(database)
}

func autoMigration(err error, db *gorm.DB) error {
	db.Migrator().DropTable(&entity.Folder{}, &entity.File{})
	err = db.AutoMigrate(&entity.Folder{}, &entity.File{})
	if err != nil {
		return err
	}
	return nil
}

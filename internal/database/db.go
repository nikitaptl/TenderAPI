package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("POSTGRES_CONN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	DB = db
	return DB, nil
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func Update(model interface{}, updates map[string]interface{}) error {
	return DB.Model(model).Updates(updates).Error
}

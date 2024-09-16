package migrations

import (
	"avitoTestTask/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, Migrations())
	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}

func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202309090000_create_uuid_extension",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec(`DROP EXTENSION IF EXISTS "uuid-ossp";`).Error
			},
		},
		{
			ID: "202309090001_create_organization_type",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`
					DO $$ 
					BEGIN
						IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
							EXECUTE 'CREATE TYPE organization_type AS ENUM (''IE'', ''LLC'', ''JSC'')';
						END IF;
					END $$;
                `).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec(`DROP TYPE IF EXISTS organization_type;`).Error
			},
		},
		{
			ID: "202309090002_create_employee_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`
					CREATE TABLE IF NOT EXISTS employee (
						id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
						username VARCHAR(50) UNIQUE NOT NULL,
						first_name VARCHAR(50),
						last_name VARCHAR(50),
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					);
				`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec(`DROP TABLE IF EXISTS employee;`).Error
			},
		},
		{
			ID: "202309090003_create_organization_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`
					CREATE TABLE IF NOT EXISTS organization (
						id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
						name VARCHAR(100) NOT NULL,
						description TEXT,
						type organization_type,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					);
				`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec(`DROP TABLE IF EXISTS organization;`).Error
			},
		},
		{
			ID: "202309090004_create_organization_responsible_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`
					CREATE TABLE IF NOT EXISTS organization_responsible (
						id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
						organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
						user_id UUID REFERENCES employee(id) ON DELETE CASCADE
					);
				`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec(`DROP TABLE IF EXISTS organization_responsible;`).Error
			},
		},
		{
			ID: "202309130001_auto_migrate_models",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Tender{}, &models.TenderVersion{}, &models.Bid{}, &models.BidVersion{}, &models.Feedback{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Tender{}, &models.TenderVersion{}, &models.Bid{}, &models.BidVersion{}, &models.Feedback{})
			},
		},
	}
}

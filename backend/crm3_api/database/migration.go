package database

import (
	"log"

	"crm3_api/models"

	"gorm.io/gorm"
)

// AutoMigrate データベースマイグレーションを実行
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Product{},
		&models.Store{},
		&models.Staff{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

// DropAllTables 全てのテーブルを削除（開発環境用）
func DropAllTables(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	err := db.Migrator().DropTable(
		&models.User{},
		&models.Customer{},
		&models.Product{},
		&models.Store{},
		&models.Staff{},
	)

	if err != nil {
		log.Printf("Failed to drop tables: %v", err)
		return err
	}

	log.Println("All tables dropped successfully")
	return nil
}
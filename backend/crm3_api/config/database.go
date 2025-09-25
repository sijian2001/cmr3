package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		Database: getEnv("DB_NAME", "crm3"),
		Username: getEnv("DB_USER", "crm3user"),
		Password: getEnv("DB_PASSWORD", "1234"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	// 本番環境で安全でないデフォルト値を使用することを防ぐ
	if isProduction() && isUnsafeDefault(key, defaultValue) {
		log.Fatalf("Production environment requires secure %s configuration. Please set the %s environment variable.", key, key)
	}

	return defaultValue
}

// isProduction は本番環境かどうかを判定する
func isProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == "production" || env == "prod"
}

// isUnsafeDefault は安全でないデフォルト値かどうかを判定する
func isUnsafeDefault(key, value string) bool {
	unsafeDefaults := map[string][]string{
		"DB_PASSWORD": {"1234", "password", "123456", "admin", "root"},
		"DB_USER":     {"crm3user", "admin", "root", "user"},
		"DB_NAME":     {}, // 本番では空でも許可する場合がある
		"DB_HOST":     {}, // localhostは開発でのみ使用
	}

	if unsafeValues, exists := unsafeDefaults[key]; exists {
		for _, unsafeValue := range unsafeValues {
			if value == unsafeValue {
				return true
			}
		}
	}

	return false
}

func (config *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}

func ConnectDatabase() (*gorm.DB, error) {
	config := NewDatabaseConfig()
	dsn := config.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}
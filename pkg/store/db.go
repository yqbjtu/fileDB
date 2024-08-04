package store

import (
	"fileDB/pkg/config"
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var MyDB *gorm.DB

// NewPgDB 从配置中新建 Postgres 存储
func InitDB() {
	postgresConfig := config.GetConfig().Postgres
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password='%s' sslmode=disable",
		postgresConfig.Host, postgresConfig.Port, postgresConfig.Username, postgresConfig.DBName, postgresConfig.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: postgresConfig.BatchSize,
		//Logger:          gl.Default.LogMode(gl.Info),
	})
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err != nil {
		fmt.Println("open dsn", dsn, "failed!", err)
	}

	db.AutoMigrate(&domain.CellStatus{})
	db.AutoMigrate(&domain.CellHistory{})
	db.AutoMigrate(&domain.CellGisMeta{})

	MyDB = db
}

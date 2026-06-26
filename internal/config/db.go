package config

import (
	parkingzone "spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// auto migrate
	db.AutoMigrate(
		&user.User{},
		&parkingzone.ParkingZone{},
	)

	println("Database connection successful")
	return db
}

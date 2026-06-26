package reservation

import (
	parkingzone "spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/user"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	UserID       uint                    `json:"user_id" gorm:"not null"`
	ZoneID       uint                    `json:"zone_id" gorm:"not null"`
	LicensePlate string                  `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       string                  `json:"status" gorm:"type:varchar(20);default:'active';not null"`
	User         user.User               `json:"user" gorm:"foreignKey:UserID"`
	Zone         parkingzone.ParkingZone `json:"zone" gorm:"foreignKey:ZoneID"`
}

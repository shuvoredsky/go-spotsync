package dto

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReservationResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	ZoneID       uint   `json:"zone_id"`
	LicensePlate string `json:"license_plate"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type MyReservationResponse struct {
	ID           uint     `json:"id"`
	LicensePlate string   `json:"license_plate"`
	Status       string   `json:"status"`
	Zone         ZoneInfo `json:"zone"`
	CreatedAt    string   `json:"created_at"`
}

type AdminReservationResponse struct {
	ID           uint     `json:"id"`
	LicensePlate string   `json:"license_plate"`
	Status       string   `json:"status"`
	Zone         ZoneInfo `json:"zone"`
	User         UserInfo `json:"user"`
	CreatedAt    string   `json:"created_at"`
}

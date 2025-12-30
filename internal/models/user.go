package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
	NumberID string `json:"number_id"`
	Role     string `gorm:"default:none" json:"role"`
}

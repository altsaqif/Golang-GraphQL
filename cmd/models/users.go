package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string `gorm:"not null" json:"name"`
	Email           string `gorm:"not null;unique" json:"email"`
	Password        string `gorm:"not null" json:"password"`
	PasswordConfirm string `gorm:"not null" json:"password_confirm"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// GetUserIdByEmail digunakan untuk mencari ID pengguna berdasarkan nama pengguna
func GetUserIdByEmail(db *gorm.DB, username string) (int, error) {
	var user User
	if err := db.Where("email = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}
	return user.ID, nil
}

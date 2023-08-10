package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `json:"-" gorm:"primary_key"`
	Name      string `gorm:"not null"`
	StudentID uint   `gorm:"unique"`
	Grade     string
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"gorm:"type:timestamptz"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

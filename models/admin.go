package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	Password  string
	CreatedAt time.Time `gorm:"type:timestamptz"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.Password = string(hashedPassword)
	return nil
}

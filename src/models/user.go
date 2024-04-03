package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          string        `gorm:"primaryKey"`
	Name        string        `gorm:"size:50"`
	Email       string        `gorm:"size:100"`
	Password    string        `gorm:"size:255"`
	PlayerCards []PlayerCards `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
}

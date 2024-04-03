package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"size:100"`
	Description string `gorm:"type:text"`
	Type        string `gorm:"size:100"`
	Position    string `gorm:"size:100"`
	Skill       string `gorm:"size:100"`
	Strength    uint
	Quantity    uint
	Faction     string `gorm:"size:100"`
	Image       string
	PlayerCards []PlayerCards `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE"`
}

func (card *Card) BeforeCreate(tx *gorm.DB) (err error) {
	card.ID = uuid.NewString()
	return
}

type PlayerCards struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	PlayerID string
	CardID   string
}

func (playerCards *PlayerCards) BeforeCreate(tx *gorm.DB) (err error) {
	playerCards.ID = uuid.NewString()
	return
}

type Match struct {
	gorm.Model
	ID            string          `gorm:"primaryKey"`
	MatchMovement []MatchMovement `gorm:"foreignKey:MatchID;constraint:OnDelete:CASCADE"`
	StartTime     time.Time
	EndTime       *time.Time `gorm:"default:null"`
}

func (match *Match) BeforeCreate(tx *gorm.DB) (err error) {
	match.ID = uuid.NewString()
	return
}

type MatchMovement struct {
	gorm.Model
	ID      string `gorm:"primaryKey"`
	MatchID string
}

func (match *MatchMovement) BeforeCreate(tx *gorm.DB) (err error) {
	match.ID = uuid.NewString()
	return
}

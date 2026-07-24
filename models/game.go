package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name         string       `gorm:"size:255;not null" json:"name" binding:"required,max=255"`
	StartedDate  *time.Time   `json:"started_date,omitempty"`
	FinishedDate *time.Time   `json:"finished_date,omitempty"`
	Played       bool         `gorm:"default:false" json:"played"`
	Description  string       `gorm:"type:text" json:"description"`
	DownloadOnly bool         `gorm:"default:false" json:"download_only"`
	Rating       int32        `json:"rating"`
	ReleaseDate  time.Time    `json:"release_date"`
	Gamesystems  []Gamesystem `gorm:"many2many:game_gamesystems;" json:"gamesystems"`
}

type GameInput struct {
	Game
	GamesystemIDs []uint `json:"gamesystem_ids" binding:"required"`
}

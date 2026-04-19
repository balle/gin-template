package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	CreatedDate  time.Time `gorm:"autoCreateTime"`
	StartedDate  *time.Time
	FinishedDate *time.Time
	Played       bool   `gorm:"default:false"`
	Description  string `gorm:"type:text"`
	DownloadOnly bool   `gorm:"default:false"`
	Rating       int32
	ReleaseDate  time.Time
}

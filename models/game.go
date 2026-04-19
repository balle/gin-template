package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name         string `gorm:"size:255;not null"`
	StartedDate  *time.Time
	FinishedDate *time.Time
	Played       bool   `gorm:"default:false"`
	Description  string `gorm:"type:text"`
	DownloadOnly bool   `gorm:"default:false"`
	Rating       int32
	ReleaseDate  time.Time
}

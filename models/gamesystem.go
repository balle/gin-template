package models

import "gorm.io/gorm"

type Gamesystem struct {
	gorm.Model
	Name string `gorm:"size:255;not null"`
}

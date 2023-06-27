package model

import (
	"time"
)

type MapInfo struct {
	MapId    string `gorm:"primaryKey"`
	Rid      string
	Nn       string
	Uid      string
	Txt      string
	Level    string
	State    string    `gorm:"default:0"`
	Star     string    `gorm:"default:0"`
	Tst      bool      `gorm:"not null"`
	PlayTime time.Time `gorm:"default:null"`
	Created  time.Time
}

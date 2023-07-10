package model

import (
	"time"
)

type Blacklist struct {
	Uid     string `gorm:"primaryKey"`
	Nn      string
	Created time.Time
}

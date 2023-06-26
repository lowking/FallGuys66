package model

import "time"

type MapInfo struct {
	MapId   string
	Rid     string
	Nn      string
	Uid     string
	Txt     string
	Level   string
	State   string
	Star    string
	Created time.Time
}

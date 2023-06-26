package handler

import (
	"FallGuys66/db"
	"FallGuys66/live/douyu/client"
	"FallGuys66/live/douyu/lib/logger"
	"regexp"
)

var mapRe = regexp.MustCompile(`([\d]{4}-[\d]{4}-[\d]{4})`)

func FilterMap(msg client.Item) {
	logger.Debugf("%s[%s][%s]: %s", msg.Nn, msg.Level, msg.Uid, msg.Txt)
	for _, mapId := range mapRe.FindAllString(msg.Txt, -1) {
		db.InsertMap(mapId, msg)
	}
}

func main() {
	FilterMap(client.Item{
		Rid:     "",
		Cid:     "",
		Uid:     "",
		Type:    "",
		Txt:     "ceshi2345-2344-36343458-3455-3455",
		Nn:      "",
		Level:   "",
		Payload: nil,
	})
}

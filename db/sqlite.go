package db

import (
	"FallGuys66/db/model"
	"FallGuys66/live/douyu/client"
	"FallGuys66/live/douyu/lib/logger"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"runtime/debug"
	"time"
)

var Db *sql.DB

func init() {
	mapDataFilePath := "./map.db"
	// 打开数据库，如果不存在，则创建
	db, err := sql.Open("sqlite3", mapDataFilePath)
	checkErr(err)
	Db = db

	// 创建表
	mapInfoTable := `
    CREATE TABLE IF NOT EXISTS mapinfo(
        mapId VARCHAR(14) PRIMARY KEY,
        rid VARCHAR(20) NULL,
        nn VARCHAR(128) NULL,
        uid VARCHAR(20) NULL,
        txt VARCHAR(1024) NULL,
        level VARCHAR(4) NULL,
        state VARCHAR(1) NOT NULL DEFAULT 0,
        star VARCHAR(1) NOT NULL DEFAULT 0,
        created DATE NULL
    );
    `

	_, err = db.Exec(mapInfoTable)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		debug.PrintStack()
		logger.Errorf("error: %v", err)
	}
}

func InsertMap(mapId string, msg client.Item) {
	// insert
	stmt, err := Db.Prepare("INSERT INTO mapinfo(mapId, rid, nn, uid, txt, level, created) values(?,?,?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(mapId, msg.Rid, msg.Nn, msg.Uid, msg.Txt, msg.Level, time.Now())
	checkErr(err)
	//
	// id, err := res.LastInsertId()
	// checkErr(err)
}

func ListMap(pageNo int, pageSize int, where string, order string) []model.MapInfo {
	start := (pageNo - 1) * pageSize
	// end := pageNo * pageSize
	rows, err := Db.Query(fmt.Sprintf("SELECT mapId, rid, nn, uid, txt, level, state, star, created FROM mapinfo where 1=1 %s %s limit %d, %d", where, order, start, pageSize))
	checkErr(err)
	var mapInfo model.MapInfo
	var mapInfos []model.MapInfo

	for rows.Next() {
		mapInfo = model.MapInfo{}
		_ = rows.Scan(
			&mapInfo.MapId,
			&mapInfo.Rid,
			&mapInfo.Nn,
			&mapInfo.Uid,
			&mapInfo.Txt,
			&mapInfo.Level,
			&mapInfo.State,
			&mapInfo.Star,
			&mapInfo.Created,
		)
		mapInfos = append(mapInfos, mapInfo)
	}
	_ = rows.Close()

	return mapInfos
}

func UpdateMap(mid string, set string, where string) {
	stmt, err := Db.Prepare(fmt.Sprintf("update mapinfo set %s where 1=1 %s", set, where))
	checkErr(err)
	_, err = stmt.Exec(mid)
	checkErr(err)
}

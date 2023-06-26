package db

import (
	"FallGuys66/live/douyu/client"
	"FallGuys66/live/douyu/lib/logger"
	"database/sql"
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
        created DATE NULL
    );
    `

	_, err = db.Exec(mapInfoTable)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		debug.Stack()
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

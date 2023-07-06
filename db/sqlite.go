package db

import (
	"FallGuys66/db/model"
	"FallGuys66/live/douyu/lib/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

var Db *gorm.DB

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	dbPath := filepath.Join(configDir, "fallguys66")
	err = os.MkdirAll(dbPath, 0755)
	if err != nil {
		panic(err)
	}
	dbPath = filepath.Join(dbPath, "map.db")
	logger.Infof("数据库文件路径：%s", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: log.Default.LogMode(log.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	Db = db
	err = Db.AutoMigrate(&model.MapInfo{})
	if err != nil {
		panic("failed to create table")
	}
}

func checkErr(err error) {
	if err != nil {
		debug.PrintStack()
		logger.Errorf("error: %v", err)
	}
}

func InsertMap(mapInfo model.MapInfo) {
	Db.Debug().Create(&mapInfo)
}

func ListMap(pageNo int, pageSize int, where *model.MapInfo, order string) []model.MapInfo {
	start := (pageNo - 1) * pageSize
	// end := pageNo * pageSize
	var mapList []model.MapInfo
	Db.Debug().Where(&where).Limit(pageSize).Offset(start).Order(order).Find(&mapList)

	return mapList
}

func SearchMap(pageNo int, pageSize int, where string, order string) []model.MapInfo {
	start := (pageNo - 1) * pageSize
	// end := pageNo * pageSize
	var mapList []model.MapInfo
	Db.Debug().Where(where).Limit(pageSize).Offset(start).Order(order).Find(&mapList)

	return mapList
}

func UpdateMap(mapInfo model.MapInfo, set []string, where *model.MapInfo) {
	Db.Debug().Where(&where).Model(&mapInfo).Select(set).Updates(&mapInfo)
}

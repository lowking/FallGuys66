package db

import (
	"FallGuys66/config"
	"FallGuys66/db/model"
	"FallGuys66/live/douyu/lib/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
	"path/filepath"
	"runtime/debug"
	"time"
)

var Db *gorm.DB

func init() {
	dbPath := filepath.Join(config.UserConfigDir, "map.db")
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
		panic(err)
	}
	err = Db.AutoMigrate(&model.Blacklist{})
	if err != nil {
		panic(err)
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

func ListMap(pageNo int, pageSize int, where *model.MapInfo, order string, excludeBlack bool) ([]model.MapInfo, int64) {
	start := (pageNo - 1) * pageSize
	// end := pageNo * pageSize
	var mapList []model.MapInfo
	var count int64
	if excludeBlack {
		Db.Debug().Where(&where).Where("uid not in (?)", Db.Model(&model.Blacklist{}).Select("uid")).Limit(pageSize).Offset(start).Order(order).Find(&mapList)
		Db.Debug().Model(&model.MapInfo{}).Where(&where).Where("uid not in (?)", Db.Model(&model.Blacklist{}).Select("uid")).Count(&count)
	} else {
		Db.Debug().Where(&where).Limit(pageSize).Offset(start).Order(order).Find(&mapList)
		Db.Debug().Model(&model.MapInfo{}).Where(&where).Count(&count)
	}

	return mapList, count
}

func SearchMap(pageNo int, pageSize int, where string, order string) ([]model.MapInfo, int64) {
	start := (pageNo - 1) * pageSize
	// end := pageNo * pageSize
	var mapList []model.MapInfo
	var count int64
	Db.Debug().Where(where).Limit(pageSize).Offset(start).Order(order).Find(&mapList)
	Db.Debug().Model(&model.MapInfo{}).Where(where).Count(&count)

	return mapList, count
}

func UpdateMap(mapInfo model.MapInfo, set []string, where *model.MapInfo) {
	Db.Debug().Where(&where).Model(&mapInfo).Select(set).Updates(&mapInfo)
}

func InsertBlacklist(blacklist model.Blacklist) {
	Db.Debug().Create(&blacklist)
}

func DeleteBlacklist(blacklist model.Blacklist) {
	Db.Debug().Delete(&blacklist)
}

func ListBlacklist(pageNo int, pageSize int, where *model.Blacklist, order string) ([]model.Blacklist, int64) {
	start := (pageNo - 1) * pageSize
	var blacklist []model.Blacklist
	var count int64
	Db.Debug().Where(&where).Limit(pageSize).Offset(start).Order(order).Find(&blacklist)
	Db.Debug().Model(&model.Blacklist{}).Where(&where).Count(&count)

	return blacklist, count
}

func CountBlacklist(where *model.Blacklist) int64 {
	var count int64
	Db.Debug().Model(&model.Blacklist{}).Where(&where).Count(&count)

	return count
}

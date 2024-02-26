package db

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	instance *leveldb.DB
)

func init() {
	// 打开或创建 LevelDB 数据库
	_db, err := leveldb.OpenFile("./run/db", nil)
	if err != nil {
		log.Fatal(err)
	}
	instance = _db
}

func Close() {
	instance.Close()
}

func GetDB() *leveldb.DB {
	return instance
}

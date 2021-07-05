package db

import (
	"gorm.io/gorm"
	"sync"
	"time"
)

var Conns sync.Map

//default ORM
var Orm *gorm.DB

func SetDefault(g *gorm.DB) {
	Orm = g
}

func ORM(con string) *gorm.DB {
	v, ok := Conns.Load(con);
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

func SetMaxConns(con string, MaxIdleConns, MaxOpenConns int) {
	orm := ORM(con)
	sqlDB, _ := orm.DB()
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
}



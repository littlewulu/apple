package mysql


import (
	"apple/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB, _ = gorm.Open("mysql", config.ConfigInstance.MysqlConfig)
)

func init(){
	if config.ConfigInstance.IsDebug {
		DB.LogMode(true)
	}
}
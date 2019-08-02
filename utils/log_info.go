package utils

import (
	"log"
	"os"
	"apple/config"
)

var (
	CommonLogFile = config.ConfigInstance.BaseDir + config.ConfigInstance.LogFilePath + "/common.log"
	ApplicationFile = config.ConfigInstance.BaseDir + config.ConfigInstance.LogFilePath + "/application.log"
	RedisFile = config.ConfigInstance.BaseDir + config.ConfigInstance.LogFilePath + "/redis.log"
	WechatFile = config.ConfigInstance.BaseDir + config.ConfigInstance.LogFilePath + "/wechat.log"
)

func Loginfo(prefix string, logFile string, data ...interface{}){
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil{
		return
	}
	defer file.Close()
	pre := "[" + prefix + "]"
	logger := log.New(file, pre, log.Ldate|log.Ltime)
	logger.Println(data...)
}

package main

import (
	"apple/config"
	"apple/controller"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func init(){
	rand.Seed(time.Now().Unix())

	if !config.ConfigInstance.IsDebug{
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if config.ConfigInstance.IsDebug{
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	}else{
		router.Use(gin.Recovery())
	}

	router.GET("/api/apple/get_application_access_token", controller.GetApplicationAccessTokenFunc)
	router.GET("/api/apple/get_wechat_access_token", controller.GetWechatAccessTokenFunc)


	// 管理调试
	router.GET("/api/apple/manager/signature", controller.GetSignatureFunc)
	router.POST("/api/apple/manager/gen_new_app", controller.GenNewAppFunc)

	router.Run(config.ConfigInstance.RestServerAddr)
}


func main(){


}



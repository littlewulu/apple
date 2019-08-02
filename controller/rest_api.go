package controller

import (
	"apple/config"
	"apple/controller/forms"
	"apple/database/models"
	"apple/manager"
	"apple/utils"
	"apple/wechat"
	"github.com/gin-gonic/gin"
)

// 请求应用token
// /api/apple/get_application_access_token
func GetApplicationAccessTokenFunc(c *gin.Context){
	var f forms.GetApplicationAccessTokenForm
	err := c.ShouldBind(&f)
	if err != nil{
		c.JSON(400, ErrorResponse{Code:10, Msg:"params error"})
		return
	}
	if ! manager.CheckAppIDAndSecret(f.AppID, f.AppSecret){
		c.JSON(400, ErrorResponse{Code:11, Msg:"params error"})
		return
	}
	m := manager.GetManagerInstance()
	token, expiredAt := m.GenerateAccessToken(f.AppID, f.AppSecret)
	if token == ""{
		c.JSON(400, ErrorResponse{Code:12, Msg:"gen token error"})
		return
	}
	c.JSON(200, map[string]interface{}{
		"access_token": token,
		"expired_at": expiredAt,
	})
	return
}

// 请求wechat token
// /api/apple/get_wechat_access_token
func GetWechatAccessTokenFunc(c *gin.Context){
	var f forms.GetWechatAccessTokenForm
	err := c.ShouldBind(&f)
	if err != nil{
		c.JSON(400, ErrorResponse{Code:10, Msg:"params error"})
		return
	}
	m := manager.GetManagerInstance()
	ok, _ := m.Validate(f.AccessToken, f.Signature)
	if ! ok{
		c.JSON(403, ErrorResponse{Code:1, Msg:"forbidden"})
		return
	}

	token, err := wechat.GetWechatAccessToken(f.Application)
	if err != nil{
		utils.Loginfo("get_wechat_access_token_func_fail", utils.CommonLogFile, map[string]interface{}{
			"form": f,
			"error": err,
		})
		c.JSON(400, ErrorResponse{Code:11, Msg:"get token fail: "+err.Error()})
		return
	}
	c.JSON(200, map[string]interface{}{
		"token": token.AccessToken,
		"expires_in": token.ExpiresIn,
		"expired_at": token.ExpiresAt,
	})
	return
}

// 管理调试

// 获取token签名 - 调试用
// /api/apple/manager/signature
func GetSignatureFunc(c *gin.Context){
	var f forms.GetSignatureForm
	err := c.ShouldBind(&f)
	if err != nil{
		c.JSON(400, ErrorResponse{Code:10, Msg:"params error"})
		return
	}
	if f.Password == "" || f.Password != config.ConfigInstance.DebugPassword{
		c.JSON(400, ErrorResponse{Code:11, Msg:"params error"})
		return
	}
	m := manager.GetManagerInstance()
	tokenInfo := m.GetATInfo(f.AccessToken)
	if tokenInfo == nil{
		c.JSON(400, ErrorResponse{Code:12, Msg:"token error"})
		return
	}
	signature := manager.CalcuSignature(f.AccessToken, tokenInfo.AppSecret)
	c.JSON(200, map[string]interface{}{
		"signature": signature,
	})
	return
}

// 添加新的app
// /api/apple/manager/gen_new_app
func GenNewAppFunc(c *gin.Context){
	var f forms.NewApplicationForm
	err := c.ShouldBind(&f)
	if err != nil{
		c.JSON(400, ErrorResponse{Code:10, Msg:"params error"})
		return
	}
	if f.Password == "" || f.Password != config.ConfigInstance.GodKey{
		c.JSON(400, ErrorResponse{Code:11, Msg:"params error"})
		return
	}

	app, err := models.CreateApplicationModel(&models.ApplicationIn{
		AppID: f.AppID,
		Description: f.Description,
	})
	if err != nil{
		c.JSON(400, ErrorResponse{Code:12, Msg:err.Error()})
		return
	}
	c.JSON(200, map[string]interface{}{
		"app_id": app.AppID,
		"app_secret": app.AppSecret,
	})
	return
}


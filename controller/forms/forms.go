package forms

type GetApplicationAccessTokenForm struct {
	AppID string `form:"app_id" json:"app_id" binding:"required"`
	AppSecret string `form:"app_secret" json:"app_secret" binding:"required"`
}

type GetWechatAccessTokenForm struct {
	Application string `form:"application" json:"application" binding:"required"`
	AccessToken string `form:"access_token" json:"access_token" binding:"required"`
	Signature string `form:"signature" json:"signature" binding:"required"`
}

type NewApplicationForm struct {
	AppID string `form:"app_id" json:"app_id" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type GetSignatureForm struct {
	AccessToken string `form:"access_token" json:"access_token" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}


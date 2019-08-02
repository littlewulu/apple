package manager

import "apple/database/models"

// 校验appID和appSecret合法性
func CheckAppIDAndSecret(appID, appSecret string)bool{
	aModel, err := models.GetApplicationModel(appID)
	if err != nil{
		return false
	}

	// 状态
	if aModel.Status == models.AppStatusBlock{
		return false
	}

	if aModel.AppSecret != appSecret{
		return false
	}

	return true
}



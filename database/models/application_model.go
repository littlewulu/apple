package models

import (
	"apple/database/mysql"
	"apple/utils"
	"errors"
	"fmt"
	"time"
)

const (
	AppStatusActive uint8 = 1
	AppStatusBlock uint8 = 2  // 封禁
)

// 应用配置表
type ApplicationModel struct {
	AppID string `gorm:"column:app_id; primary_key"`
	AppSecret string `gorm:"column:app_secret"`
	Description string `gorm:"column:description"`
	Status uint8 `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (a *ApplicationModel)TableName()string{
	return "application_config"
}

func (a *ApplicationModel)UpdateStatus(s uint8){
	mysql.DB.Model(a).UpdateColumn("status", s)
}

type ApplicationIn struct {
	AppID string
	Description string
}

func CreateApplicationModel(i *ApplicationIn)(*ApplicationModel, error){
	if i.AppID == ""{
		return nil, errors.New("app_id cannot be empty")
	}
	aTmp, _ := GetApplicationModel(i.AppID)
	if aTmp != nil{
		return nil, errors.New("app_id already exist")
	}

	a := ApplicationModel{
		AppID: i.AppID,
		AppSecret: genAppSecret(i.AppID),
		Description: i.Description,
		Status: AppStatusActive,
	}
	if err := mysql.DB.Create(&a).Error; err != nil{
		return nil, err
	}
	return &a, nil
}

func GetApplicationModel(aID string)(*ApplicationModel, error){
	a := ApplicationModel{}
	if err := mysql.DB.Where("app_id = ?", aID).Find(&a).Error; err != nil{
		return nil, err
	}
	return &a, nil
}


func genAppSecret(appID string)string{
	s := fmt.Sprint(time.Now().UnixNano()) + appID + utils.RanString(128)
	r := utils.Md5(s)
	return r
}





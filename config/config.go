package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	configFilePath = "/Users/sherlock/Desktop/my_project/apple/src/apple/config/config.yaml"
	ConfigInstance ConfigType
)

type WechatApp struct {
	AppID string `yaml:"AppID"`
	AppSecret string `yaml:"AppSecret"`
}

// 配置结构
type ConfigType struct {
	IsDebug bool `yaml:"IsDebug"`
	BaseDir string `yaml:"BaseDir"`
	// redis 配置
	RedisAddr string `yaml:"RedisAddr"`
	RedisPassword string `yaml:"RedisPassword"`
	RedisDatabase int `yaml:"RedisDatabase"`
	RedisServerPrefixKey string `yaml:"RedisServerPrefixKey"`
	// mysql
	MysqlConfig string `yaml:"MysqlConfig"`
	// 日志
	LogFilePath string `yaml:"LogFilePath"`
	// 调试密码
	DebugPassword string `yaml:"DebugPassword"`
	// access_token 配置
	AccessTokenEffectiveTime int64 `yaml:"AccessTokenEffectiveTime"`
	AccessTokenExtensionTime int64 `yaml:"AccessTokenExtensionTime"`
	// 监听地址
	RestServerAddr string `yaml:"RestServerAddr"`
	// 操作密钥
	GodKey string `yaml:"GodKey"`
	// wechat 配置
	WechatAppConfig map[string]WechatApp `yaml:"WechatAppConfig"`

}

func init(){
	configFileEnvPath := os.Getenv("CONFIG_FILE_PATH")
	if configFileEnvPath != ""{
		configFilePath = configFileEnvPath
	}
	fileData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(fileData, &ConfigInstance)
	if err != nil{
		panic(err)
	}
}



package wechat

import (
	"apple/config"
	"apple/database/redisclient"
	"apple/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// 小程序
	TinyProgramXueJun = "tiny_program_xue_jun"


	// wechat 接口
	wechatUrlGetAccessToken = "https://api.weixin.qq.com/cgi-bin/token"

	// redis 存储
	redisPrefix = "apple:wechat_access_token:"
	redisSub = 1800  // 少存半个钟
)


type WechatAccessToken struct {
	App string
	AccessToken string
	ExpiresIn int64
	ExpiresAt int64
}

type wechatAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
}

// 获取微信应用的access_token
func GetWechatAccessToken(app string)(*WechatAccessToken, error){
	// 先读缓存
	token, err := getAccessTokenFromRedis(app)
	if err != nil{
		utils.Loginfo("get_access_token_from_redis_error", utils.WechatFile, err)
	}
	if token != nil{
		return token, nil
	}

	// 向wechat请求
	token, err = getAccessTokenFromWechat(app)
	if err != nil{
		utils.Loginfo("get_access_token_from_wechat_error", utils.WechatFile, err)
		return nil, err
	}
	// 设置缓存
	err = setAccessTokenToRedis(token)
	if err != nil{
		utils.Loginfo("set_access_token_to_redis_error", utils.WechatFile, err)
		return nil, err
	}
	return token, nil
}

func getAccessTokenFromRedis(app string)(*WechatAccessToken, error){
	rKey := getRedisKey(app)
	data, err := redisclient.RedisClient.Get(rKey).Result()
	if err != nil{
		return nil, err
	}
	if data == ""{
		return nil, nil
	}
	token := WechatAccessToken{}
	err = json.Unmarshal([]byte(data), &token)
	if err != nil{
		return nil, err
	}
	token.ExpiresIn = token.ExpiresAt - time.Now().Unix()
	return &token, nil
}

func setAccessTokenToRedis(w *WechatAccessToken)error{
	rKey := getRedisKey(w.App)
	data, err := json.Marshal(w)
	if err != nil{
		return err
	}
	// 存少半个钟
	_, err = redisclient.RedisClient.Set(rKey, string(data), time.Second * time.Duration(w.ExpiresIn)).Result()
	return err
}

func getRedisKey(app string)string{
	return redisPrefix + config.ConfigInstance.RedisServerPrefixKey + app
}


func getAccessTokenFromWechat(app string)(*WechatAccessToken, error){
	appConfig, ok := config.ConfigInstance.WechatAppConfig[app]
	if ! ok{
		return nil, errors.New("app not exist")
	}
	params := url.Values{
		"grant_type": []string{"client_credential"},
		"appid": []string{appConfig.AppID},
		"secret": []string{appConfig.AppSecret},
	}
	api := wechatUrlGetAccessToken + "?" + params.Encode()
	res, err := http.Get(api)
	if err != nil{
		utils.Loginfo("get_access_token_from_wechat_fail", utils.WechatFile, map[string]interface{}{
			"error": err,
			"res": res,
			"api": api,
		})
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil{
		utils.Loginfo("get_access_token_from_wechat_fail", utils.WechatFile, map[string]interface{}{
			"error": err,
			"res": res,
			"api": api,
		})
		return nil, err
	}
	result := wechatAccessTokenResponse{}
	err = json.Unmarshal(resBody, &result)
	if err != nil{
		utils.Loginfo("get_access_token_from_wechat_fail", utils.WechatFile, map[string]interface{}{
			"error": err,
			"res": res,
			"api": api,
			"body": string(resBody),
		})
		return nil, err
	}

	token := WechatAccessToken{
		App: app,
		AccessToken: result.AccessToken,
		ExpiresIn: result.ExpiresIn - redisSub,
		ExpiresAt: time.Now().Unix() + result.ExpiresIn - redisSub ,
	}

	utils.Loginfo("get_access_token_from_wechat_success", utils.WechatFile, token)
	return &token, nil
}




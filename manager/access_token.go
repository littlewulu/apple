package manager

import (
	"apple/config"
	"apple/database/redisclient"
	"apple/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// 服务间调用凭证使用redis存储
// 调用凭证access_token的有效期是2小时
// 2小时过期后，该调用凭证纳入过期集合，延长5分钟的有效时间
/*
存储结构为

{
	access_token1 : {
		"app_id": "xxxx",
		"app_secret": "xxxx",
		"expired_at": timestamp
	}
}
*/

var (
	effectiveTime int64 = config.ConfigInstance.AccessTokenEffectiveTime
	extensionTime int64 = config.ConfigInstance.AccessTokenExtensionTime

	// 管理器实例
	accessTokenManager = &AccessTokenManagerStruct{}

	// redis存储集合 键值
	accessTokenSetRedisKey = "apple:access_token_set:" + config.ConfigInstance.RedisServerPrefixKey
	accessTokenExtensionSetRedisKey = "apple:access_token_extension_set:" + config.ConfigInstance.RedisServerPrefixKey

)

const (
	// token 合法性检验
	TokenStatusYes int = 1
	TokenStatusNo int = 2
)

// access_token 保存在redis中的结构
type AccessTokenRedisStruct struct {
	AppID string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	ExpiredAt int64 `json:"expired_at"`
}

// 使用一个空壳结构来内聚管理器方法
type AccessTokenManagerStruct struct {}

func (a *AccessTokenRedisStruct)IsExpired()bool{
	if a.ExpiredAt > time.Now().Unix(){
		return false
	}
	return true
}

// 获取实例
func GetManagerInstance()*AccessTokenManagerStruct{
	return accessTokenManager
}


// 读取access_token info
// 同时处理过期token
// 若token过期或者不存在，返回nil
func (s *AccessTokenManagerStruct)GetATInfo(accessToken string)*AccessTokenRedisStruct{
	// 先读正集
	isExt := false
	rawStr := redisclient.HGet(accessTokenSetRedisKey, accessToken)
	if rawStr == ""{
		// 再读过期集
		rawStr = redisclient.HGet(accessTokenExtensionSetRedisKey, accessToken)
		if rawStr == ""{
			return nil
		}
		isExt = true
	}
	a := &AccessTokenRedisStruct{}
	err := json.Unmarshal([]byte(rawStr), a)
	if err != nil{
		utils.Loginfo("get_access_token_error", utils.ApplicationFile, map[string]interface{}{
			"error": err,
			"token": accessToken,
		})
		return nil
	}

	// 判断是否过期
	if a.IsExpired(){
		if isExt{
			redisclient.HDel(accessTokenExtensionSetRedisKey, accessToken)
			return nil
		}

		redisclient.HDel(accessTokenSetRedisKey, accessToken)
		// 是否可以延长时间
		if a.ExpiredAt + extensionTime < time.Now().Unix(){
			return nil
		}else{
			// 添加到过期集合
			s.SetATToExtSet(accessToken, a, 0)
			return a
		}
	}
	return a
}


// 设置access_token
func (s *AccessTokenManagerStruct)SetAT(accessToken string, ar *AccessTokenRedisStruct)error{
	raw, err := json.Marshal(ar)
	if err != nil{
		utils.Loginfo("set_access_token_error", utils.ApplicationFile, map[string]interface{}{
			"error": err,
			"token": accessToken,
		})
		return err
	}
	err = redisclient.HSet(accessTokenSetRedisKey, accessToken, raw)
	if err != nil{
		return err
	}
	// 异步清理过期token
	if rand.Intn(100) < 10{
		go func() {
			s.ClearExpired()
		}()
	}
	return nil
}

// 设置access_token 到过期集合
func (s *AccessTokenManagerStruct)SetATToExtSet(accessToken string, ar *AccessTokenRedisStruct, now int64)error{
	// 过期时间延长5分钟
	if now == 0{
		ar.ExpiredAt += extensionTime
	}else{
		ar.ExpiredAt = now + extensionTime
	}
	raw, err := json.Marshal(ar)
	if err != nil{
		utils.Loginfo("set_ext_access_token_error", utils.ApplicationFile, map[string]interface{}{
			"error": err,
			"token": accessToken,
		})
		return err
	}
	err = redisclient.HSet(accessTokenExtensionSetRedisKey, accessToken, raw)
	if err != nil{
		utils.Loginfo("set_ext_access_token_error", utils.ApplicationFile, map[string]interface{}{
			"error": err,
			"token": accessToken,
		})
		return err
	}

	// 异步清理过期token
	if rand.Intn(100) < 10{
		go func() {
			s.ClearExtExpired()
		}()
	}
	return nil
}

// 删除access_token
func (s *AccessTokenManagerStruct)RmAT(accessToken string){
	redisclient.HDel(accessTokenSetRedisKey, accessToken)
}

// 过期集合中删除access_token
func (s *AccessTokenManagerStruct)RmExtAT(accessToken string){
	redisclient.HDel(accessTokenExtensionSetRedisKey, accessToken)
}

// 清理过期的access_token
func (s *AccessTokenManagerStruct)ClearExpired(){
	rMap := redisclient.HGetAll(accessTokenSetRedisKey)
	if rMap == nil{
		return
	}

	for token := range rMap{
		rawStr := rMap[token]
		ar := &AccessTokenRedisStruct{}
		err := json.Unmarshal([]byte(rawStr), ar)
		if err != nil{
			redisclient.HDel(accessTokenSetRedisKey, token)
			continue
		}
		// 过期了
		if ar.ExpiredAt < time.Now().Unix(){
			// 丢进过期集合
			s.SetATToExtSet(token, ar, 0)
			// 删除
			redisclient.HDel(accessTokenSetRedisKey, token)
		}
	}
}

// 清理过期集合
func (s *AccessTokenManagerStruct)ClearExtExpired(){
	rMap := redisclient.HGetAll(accessTokenExtensionSetRedisKey)
	if rMap == nil{
		return
	}
	for token := range rMap{
		rawStr := rMap[token]
		ar := &AccessTokenRedisStruct{}
		err := json.Unmarshal([]byte(rawStr), ar)
		if err != nil{
			redisclient.HDel(accessTokenExtensionSetRedisKey, token)
			continue
		}
		if ar.ExpiredAt < time.Now().Unix(){
			redisclient.HDel(accessTokenExtensionSetRedisKey, token)
		}
	}
}

func (s *AccessTokenManagerStruct)GetTokenByAppID(appID string)(string, *AccessTokenRedisStruct){
	rMap := redisclient.HGetAll(accessTokenSetRedisKey)
	if rMap == nil{
		return "", nil
	}

	for token := range rMap {
		rawStr := rMap[token]
		ar := &AccessTokenRedisStruct{}
		err := json.Unmarshal([]byte(rawStr), ar)
		if err != nil {
			redisclient.HDel(accessTokenSetRedisKey, token)
			continue
		}
		if ar.AppID == appID{
			return token, ar
		}
	}
	return "", nil
}

func (s *AccessTokenManagerStruct)GenerateAccessToken(appID, appSecret string)(string, int64){
	// 先查找appID是否有存在的token
	tmpToken, tmpAr := s.GetTokenByAppID(appID)
	if tmpToken != ""{
		// 失效掉原token
		s.RmAT(tmpToken)
		// 丢进过期集合
		now := time.Now().Unix()
		if tmpAr.IsExpired(){
			now = 0
		}
		s.SetATToExtSet(tmpToken, tmpAr, now)
	}

	// 生成新的token
	source := utils.RanString(32) + appID + appSecret + fmt.Sprint(time.Now().Unix())
	token := utils.Md5(source)
	ar := &AccessTokenRedisStruct{
		AppID: appID,
		AppSecret: appSecret,
		ExpiredAt: time.Now().Unix() + effectiveTime,
	}
	err := s.SetAT(token, ar)
	if err != nil{
		utils.Loginfo("generate_access_token_fail", utils.ApplicationFile, map[string]interface{}{
			"error": err,
			"token": ar,
		})
		return "", 0
	}
	return token, ar.ExpiredAt
}

// 检验token合法性
func (s *AccessTokenManagerStruct)Validate(accessToken, signature string)(bool, string){
	// token 是否存在
	info := s.GetATInfo(accessToken)
	if info == nil{
		return false, ""
	}

	// 校验签名
	tmpSig := CalcuSignature(accessToken, info.AppSecret)
	if tmpSig != signature{
		return false, ""
	}
	return true, info.AppID
}



